// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs -- -I ../../ebpf/c -I ../../../ebpf/c -fsigned-char types.go

package http2

const (
	maxHTTP2Path     = 0xa0
	http2PathBuckets = 0x7

	HTTP2TerminatedBatchSize = 0x50

	http2RawStatusCodeMaxLength = 0x3
)

type connTuple = struct {
	Saddr_h  uint64
	Saddr_l  uint64
	Daddr_h  uint64
	Daddr_l  uint64
	Sport    uint16
	Dport    uint16
	Netns    uint32
	Pid      uint32
	Metadata uint32
}
type HTTP2DynamicTableIndex struct {
	Index uint64
	Tup   connTuple
}
type HTTP2DynamicTableEntry struct {
	Buffer             [160]int8
	Original_index     uint32
	String_len         uint8
	Is_huffman_encoded bool
	Pad_cgo_0          [2]byte
}
type http2StreamKey struct {
	Tup       connTuple
	Id        uint32
	Pad_cgo_0 [4]byte
}
type http2StatusCode struct {
	Raw_buffer         [3]uint8
	Is_huffman_encoded bool
	Indexed_value      uint8
	Finalized          bool
}
type http2Stream struct {
	Response_last_seen    uint64
	Request_started       uint64
	Status_code           http2StatusCode
	Request_method        uint8
	Path_size             uint8
	Request_end_of_stream bool
	Is_huffman_encoded    bool
	Pad_cgo_0             [6]byte
	Request_path          [160]uint8
}
type EbpfTx struct {
	Tuple  connTuple
	Stream http2Stream
}
type HTTP2Telemetry struct {
	Request_seen                     uint64
	Response_seen                    uint64
	End_of_stream                    uint64
	End_of_stream_rst                uint64
	Path_exceeds_frame               uint64
	Exceeding_max_interesting_frames uint64
	Exceeding_max_frames_to_filter   uint64
	Path_size_bucket                 [8]uint64
}

type StaticTableEnumValue = uint8

const (
	GetValue       StaticTableEnumValue = 0x2
	PostValue      StaticTableEnumValue = 0x3
	EmptyPathValue StaticTableEnumValue = 0x4
	IndexPathValue StaticTableEnumValue = 0x5
	K200Value      StaticTableEnumValue = 0x8
	K204Value      StaticTableEnumValue = 0x9
	K206Value      StaticTableEnumValue = 0xa
	K304Value      StaticTableEnumValue = 0xb
	K400Value      StaticTableEnumValue = 0xc
	K404Value      StaticTableEnumValue = 0xd
	K500Value      StaticTableEnumValue = 0xe
)
