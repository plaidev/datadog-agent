package flowaggregator

import (
	"sync"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/log"

	"github.com/DataDog/datadog-agent/pkg/netflow/common"
)

const flowFlushInterval = 60 // TODO: make it configurable
const flowContextTTL = flowFlushInterval * 5

var timeNow = time.Now

// floWrapper contains flow information and additional flush related data
type flowWrapper struct {
	flow                *common.Flow
	nextFlush           time.Time
	lastSuccessfulFlush time.Time
}

// flowAccumulator is used to accumulate aggregated flows
type flowAccumulator struct {
	flows map[string]flowWrapper
	mu    sync.Mutex
}

func newFlowWrapper(flow *common.Flow) flowWrapper {
	now := timeNow()
	return flowWrapper{
		flow:      flow,
		nextFlush: now,
	}
}

func newFlowAccumulator() *flowAccumulator {
	return &flowAccumulator{
		flows: make(map[string]flowWrapper),
	}
}

func (f *flowAccumulator) flush() []*common.Flow {
	f.mu.Lock()
	defer f.mu.Unlock()

	var flows []*common.Flow
	for key, flow := range f.flows {
		now := timeNow()
		if flow.nextFlush.After(now) {
			continue
		}
		if flow.flow != nil {
			flows = append(flows, flow.flow)
			flow.lastSuccessfulFlush = now
			flow.flow = nil
		} else if flow.lastSuccessfulFlush.Add(flowContextTTL * time.Second).Before(now) {
			// delete flow wrapper if there is no successful flushes since `flowContextTTL`
			delete(f.flows, key)
			continue
		}
		flow.nextFlush = flow.nextFlush.Add(flowFlushInterval * time.Second)
		f.flows[key] = flow
	}
	return flows
}

func (f *flowAccumulator) add(flowToAdd *common.Flow) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// TODO: handle port direction (see network-http-logger)
	// TODO: ignore ephemeral ports

	log.Tracef("New Flow (digest=%s): %+v", flowToAdd.AggregationHash(), flowToAdd)

	aggFlow, ok := f.flows[flowToAdd.AggregationHash()]
	aggHash := flowToAdd.AggregationHash()
	if !ok {
		f.flows[aggHash] = newFlowWrapper(flowToAdd)
	} else {
		if aggFlow.flow == nil {
			aggFlow.flow = flowToAdd
		} else {
			aggFlow.flow.Bytes += flowToAdd.Bytes
			aggFlow.flow.Packets += flowToAdd.Packets
			aggFlow.flow.ReceivedTimestamp = common.MinUint64(aggFlow.flow.ReceivedTimestamp, flowToAdd.ReceivedTimestamp)
			aggFlow.flow.StartTimestamp = common.MinUint64(aggFlow.flow.StartTimestamp, flowToAdd.StartTimestamp)
			aggFlow.flow.EndTimestamp = common.MaxUint64(aggFlow.flow.EndTimestamp, flowToAdd.EndTimestamp)

			// TODO: Cumulate TCPFlags (Cumulative of all the TCP flags seen for this flow)
		}
		f.flows[aggHash] = aggFlow
	}
}
