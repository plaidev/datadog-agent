// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package payload

// Exporter contains exporter details
type Exporter struct {
	IP string `json:"ip"`
}

// Endpoint contains source or destination endpoint details
type Endpoint struct {
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
	// TODO: mac address
	// TODO: mask
}

// Interface contains interface details
type Interface struct {
	Index uint32 `json:"index"`
}

// ObservationPoint contains ingress or egress observation point
type ObservationPoint struct {
	Interface Interface `json:"interface"`
}

// FlowPayload contains network devices flows
type FlowPayload struct {
	FlowType string `json:"type"`
	//Timestamp    uint64           `json:"timestamp"`
	SamplingRate uint64           `json:"sampling_rate"`
	Direction    string           `json:"direction"`
	Start        uint64           `json:"start"`
	End          uint64           `json:"end"`
	Bytes        uint64           `json:"bytes"`
	Packets      uint64           `json:"packets"`
	EtherType    uint32           `json:"ether_type"`
	IPProtocol   uint32           `json:"ip_protocol"`
	Tos          uint32           `json:"tos"`
	Exporter     Exporter         `json:"exporter"`
	Source       Endpoint         `json:"flow_source"`
	Destination  Endpoint         `json:"flow_destination"`
	Ingress      ObservationPoint `json:"ingress"`
	Egress       ObservationPoint `json:"egress"`
	// TODO: Tags
	// TODO: tcp_flags
	// TODO: next_hop IP
}
