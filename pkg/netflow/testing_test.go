package netflow

import (
	"github.com/DataDog/datadog-agent/pkg/netflow/flowaggregator"
	"time"
)

func waitFlowsToBeFlushed(flowAgg *flowaggregator.FlowAggregator, timeout time.Duration) int {
	ticker := time.NewTicker(10 * time.Millisecond)
	timeoutOn := time.Now().Add(timeout)
	for {
		select {
		case <-ticker.C:
			flushCount := flowAgg.Flush()

			// this case could always take priority on the timeout case, we have to make sure
			// we've not timeout
			if time.Now().After(timeoutOn) {
				return flushCount
			}

			if flushCount > 0 {
				return flushCount
			}
		case <-time.After(timeout):
			return 0
		}
	}
}
