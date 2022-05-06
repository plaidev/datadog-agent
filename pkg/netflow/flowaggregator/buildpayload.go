package flowaggregator

import (
	"context"
	coreconfig "github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"

	coreutil "github.com/DataDog/datadog-agent/pkg/util"

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

	hostname, err := coreutil.GetHostname(context.TODO())
	if err != nil {
		log.Warnf("Error getting the hostname: %v", err)
	}

	namespace := coreconfig.Datadog.GetString("network_devices.namespace")

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
			// TODO: implement Mac
			// TODO: implement Mask
			Mac:  "00:00:00:00:00:00",
			Mask: "0.0.0.0/24",
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
		Namespace: namespace,
		Host:      hostname,
		// TODO: implement tcp_flags
		TcpFlags: []string{"SYN", "ACK"},
		NextHop: payload.NextHop{
			IP: aggFlow.NextHop,
		},
	}
}
