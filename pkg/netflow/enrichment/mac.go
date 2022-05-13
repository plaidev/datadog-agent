package enrichment

import (
	"encoding/binary"
	"net"
)

func FormatMacAddress(fieldValue uint64) string {
	mac := make([]byte, 8)
	binary.BigEndian.PutUint64(mac, fieldValue)
	return net.HardwareAddr(mac[2:]).String()
}
