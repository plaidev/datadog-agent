package flowaggregator

import (
	"github.com/DataDog/datadog-agent/pkg/netflow/common"
	"github.com/DataDog/datadog-agent/pkg/netflow/payload"
)

func buildPayload(aggFlow *common.Flow) payload.FlowPayload {
	var direction string
	if aggFlow.Direction == 0 {
		direction = "ingress"
	} else {
		direction = "egress"
	}
	return payload.FlowPayload{
		FlowType: string(aggFlow.FlowType),
		//Timestamp:    aggFlow.ReceivedTimestamp,
		SamplingRate: aggFlow.SamplingRate,
		Direction:    direction,
		Exporter: payload.Exporter{
			IP: aggFlow.SamplerAddr,
		},
		Start:      aggFlow.StartTimestamp,
		End:        aggFlow.EndTimestamp,
		Bytes:      aggFlow.Bytes,
		Packets:    aggFlow.Packets,
		EtherType:  aggFlow.EtherType,
		IPProtocol: aggFlow.IPProtocol,
		Tos:        aggFlow.Tos,
		Source: payload.Endpoint{
			IP:   aggFlow.SrcAddr,
			Port: aggFlow.SrcPort,
		},
		Destination: payload.Endpoint{
			IP:   aggFlow.DstAddr,
			Port: aggFlow.DstPort,
		},
		Ingress: payload.ObservationPoint{
			Interface: payload.Interface{
				Index: aggFlow.InputInterface,
			},
		},
		Egress: payload.ObservationPoint{
			Interface: payload.Interface{
				Index: aggFlow.OutputInterface,
			},
		},
	}
}
