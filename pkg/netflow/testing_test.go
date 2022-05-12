package netflow

import (
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/netflow/flowaggregator"
	"github.com/netsampler/goflow2/utils"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"time"
)

type dummyFlowProcessor struct {
	receivedMessages chan interface{}
	stopped          bool
}

func (d *dummyFlowProcessor) FlowRoutine(workers int, addr string, port int, reuseport bool) error {
	return utils.UDPStoppableRoutine(make(chan struct{}), "test_udp", func(msg interface{}) error {
		d.receivedMessages <- msg
		return nil
	}, 3, addr, port, false, logrus.StandardLogger())
}

func (d *dummyFlowProcessor) Shutdown() {
	d.stopped = true
}

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

func getFreePort() uint16 {
	var port uint16
	for i := 0; i < 5; i++ {
		conn, err := net.ListenPacket("udp", ":0")
		if err != nil {
			continue
		}
		conn.Close()
		port, err = parsePort(conn.LocalAddr().String())
		if err != nil {
			continue
		}
		return port
	}
	panic("unable to find free port for starting the trap listener")
}

func parsePort(addr string) (uint16, error) {
	_, portString, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}

	port, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(port), nil
}

func sendUDPPacket(port uint16, data []byte) error {
	udpConn, err := net.Dial("udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}
	_, err = udpConn.Write(data)
	udpConn.Close()
	return err
}
