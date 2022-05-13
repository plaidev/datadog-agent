package common

import "testing"

func TestFlow_AggregationHash(t *testing.T) {
	type fields struct {
		FlowType          FlowType
		ReceivedTimestamp uint64
		SamplingRate      uint64
		Direction         uint32
		SamplerAddr       string
		StartTimestamp    uint64
		EndTimestamp      uint64
		Bytes             uint64
		Packets           uint64
		SrcAddr           string
		DstAddr           string
		EtherType         uint32
		IPProtocol        uint32
		SrcPort           uint32
		DstPort           uint32
		InputInterface    uint32
		OutputInterface   uint32
		SrcMac            uint64
		DstMac            uint64
		SrcMask           uint32
		DstMask           uint32
		Tos               uint32
		NextHop           string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Flow{
				FlowType:          tt.fields.FlowType,
				ReceivedTimestamp: tt.fields.ReceivedTimestamp,
				SamplingRate:      tt.fields.SamplingRate,
				Direction:         tt.fields.Direction,
				SamplerAddr:       tt.fields.SamplerAddr,
				StartTimestamp:    tt.fields.StartTimestamp,
				EndTimestamp:      tt.fields.EndTimestamp,
				Bytes:             tt.fields.Bytes,
				Packets:           tt.fields.Packets,
				SrcAddr:           tt.fields.SrcAddr,
				DstAddr:           tt.fields.DstAddr,
				EtherType:         tt.fields.EtherType,
				IPProtocol:        tt.fields.IPProtocol,
				SrcPort:           tt.fields.SrcPort,
				DstPort:           tt.fields.DstPort,
				InputInterface:    tt.fields.InputInterface,
				OutputInterface:   tt.fields.OutputInterface,
				SrcMac:            tt.fields.SrcMac,
				DstMac:            tt.fields.DstMac,
				SrcMask:           tt.fields.SrcMask,
				DstMask:           tt.fields.DstMask,
				Tos:               tt.fields.Tos,
				NextHop:           tt.fields.NextHop,
			}
			if got := f.AggregationHash(); got != tt.want {
				t.Errorf("AggregationHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
