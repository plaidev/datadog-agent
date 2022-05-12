package netflow

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/aggregator/mocksender"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestNewNetflowServer(t *testing.T) {
	// NetFlow5 example data from goflow repo:
	// https://github.com/netsampler/goflow2/blob/5300494e478567a0fb9b8de0c504d9442260d0ad/decoders/netflowlegacy/netflow_test.go#L11-L32
	netflowV5Data := []byte{
		0x00, 0x05, 0x00, 0x06, 0x00, 0x82, 0xc3, 0x48, 0x5b, 0xcd, 0xba, 0x1b, 0x05, 0x97, 0x6d, 0xc7,
		0x00, 0x00, 0x64, 0x3d, 0x08, 0x08, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x79, 0x0a, 0x80, 0x02, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x02, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x02, 0x4e,
		0x00, 0x82, 0x9b, 0x8c, 0x00, 0x82, 0x9b, 0x90, 0x1f, 0x90, 0xb9, 0x18, 0x00, 0x1b, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x77, 0x0a, 0x81, 0x02, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x94,
		0x00, 0x82, 0x95, 0xa9, 0x00, 0x82, 0x9a, 0xfb, 0x1f, 0x90, 0xc1, 0x2c, 0x00, 0x12, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x81, 0x02, 0x01, 0x0a, 0x80, 0x02, 0x77,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x07, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0xc2,
		0x00, 0x82, 0x95, 0xa9, 0x00, 0x82, 0x9a, 0xfc, 0xc1, 0x2c, 0x1f, 0x90, 0x00, 0x16, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x01, 0x0a, 0x80, 0x02, 0x79,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x09, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x01, 0xf1,
		0x00, 0x82, 0x9b, 0x8c, 0x00, 0x82, 0x9b, 0x8f, 0xb9, 0x18, 0x1f, 0x90, 0x00, 0x1b, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x01, 0x0a, 0x80, 0x02, 0x79,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x09, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x02, 0x2e,
		0x00, 0x82, 0x9b, 0x90, 0x00, 0x82, 0x9b, 0x9d, 0xb9, 0x1a, 0x1f, 0x90, 0x00, 0x1b, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x80, 0x02, 0x79, 0x0a, 0x80, 0x02, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x02, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x0b, 0xac,
		0x00, 0x82, 0x9b, 0x90, 0x00, 0x82, 0x9b, 0x9d, 0x1f, 0x90, 0xb9, 0x1a, 0x00, 0x1b, 0x06, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	// Setup NetFlow feature config
	port := getFreePort()
	config.Datadog.SetConfigType("yaml")
	err := config.Datadog.MergeConfigOverride(strings.NewReader(fmt.Sprintf(`
network_devices:
  netflow:
    enabled: true
    aggregator_flush_interval: 1
    listeners:
      - flow_type: netflow5 # netflow, sflow, ipfix
        bind_host: 0.0.0.0
        port: %d # default 2055 for netflow
`, port)))
	require.NoError(t, err)

	// Setup NetFlow Server
	demux := aggregator.InitTestAgentDemultiplexerWithFlushInterval(10 * time.Millisecond)
	defer demux.Stop(false)
	sender := mocksender.NewMockSender("")
	sender.On("EventPlatformEvent", mock.Anything, mock.Anything).Return()
	sender.On("Count", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	demux.SetMockDefaultSender(sender)

	server, err := NewNetflowServer(demux)
	require.NoError(t, err, "cannot start Netflow Server")
	assert.NotNil(t, server)
	defer server.stop()

	// Send netflowV5Data twice to test aggregator
	// Flows will have 2x bytes/packets after aggregation
	err = sendUDPPacket(port, netflowV5Data)
	require.NoError(t, err)
	err = sendUDPPacket(port, netflowV5Data)
	require.NoError(t, err)

	// Make sure flows are flushed before doing assertions
	waitFlowsToBeFlushed(server.flowAgg, 10*time.Second)

	// Assertions

	// language = json
	event := []byte(`
{
  "type": "netflow5",
  "sampling_rate": 0,
  "direction": "ingress",
  "start": 1540209168,
  "end": 1540209169,
  "bytes": 388,
  "packets": 6,
  "ether_type": 2048,
  "ip_protocol": 6,
  "tos": 0,
  "exporter": {
    "ip": "127.0.0.1"
  },
  "source": {
    "ip": "10.129.2.1",
    "port": 49452,
    "mac": "00:00:00:00:00:00",
    "mask": "0.0.0.0/24"
  },
  "destination": {
    "ip": "10.128.2.119",
    "port": 8080,
    "mac": "",
    "mask": ""
  },
  "ingress": {
    "interface": {
      "index": 1
    }
  },
  "egress": {
    "interface": {
      "index": 7
    }
  },
  "namespace": "default",
  "host": "COMP-C02CF0CWLVDP",
  "tcp_flags": [
    "SYN",
    "ACK"
  ],
  "next_hop": {
    "ip": "0.0.0.0"
  }
}
`)
	compactEvent := new(bytes.Buffer)
	err = json.Compact(compactEvent, event)
	assert.NoError(t, err)

	sender.AssertEventPlatformEvent(t, compactEvent.String(), "network-devices-netflow")
	sender.AssertMetric(t, "Count", "datadog.newflow.aggregator.flows_received", 1, "", []string{"sample_addr:127.0.0.1", "flow_type:netflow5"})
}

func TestStartServerAndStopServer(t *testing.T) {
	demux := aggregator.InitTestAgentDemultiplexerWithFlushInterval(10 * time.Millisecond)
	defer demux.Stop(false)
	err := StartServer(demux)
	require.NoError(t, err)
	require.NotNil(t, serverInstance)

	StopServer()
	require.Nil(t, serverInstance)
}

func TestIsEnabled(t *testing.T) {
	saved := config.Datadog.Get("network_devices.netflow.enabled")
	defer config.Datadog.Set("network_devices.netflow.enabled", saved)

	config.Datadog.Set("network_devices.netflow.enabled", true)
	assert.Equal(t, true, IsEnabled())

	config.Datadog.Set("network_devices.netflow.enabled", false)
	assert.Equal(t, false, IsEnabled())
}

func TestServer_Stop(t *testing.T) {
	// Setup logger to record logs
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	l, err := seelog.LoggerFromWriterWithMinLevelAndFormat(w, seelog.DebugLvl, "[%LEVEL] %FuncShort: %Msg")
	assert.Nil(t, err)
	log.SetupLogger(l, "debug")

	// Setup NetFlow config
	port := getFreePort()
	config.Datadog.SetConfigType("yaml")
	err = config.Datadog.MergeConfigOverride(strings.NewReader(fmt.Sprintf(`
network_devices:
  netflow:
    enabled: true
    aggregator_flush_interval: 1
    listeners:
      - flow_type: netflow5 # netflow, sflow, ipfix
        bind_host: 0.0.0.0
        port: %d # default 2055 for netflow
`, port)))
	require.NoError(t, err)

	// Setup Netflow Server
	demux := aggregator.InitTestAgentDemultiplexerWithFlushInterval(10 * time.Millisecond)
	defer demux.Stop(false)
	server, err := NewNetflowServer(demux)
	require.NoError(t, err, "cannot start Netflow Server")
	assert.NotNil(t, server)

	// Stops server
	server.stop()

	// Assert logs present
	w.Flush()
	logs := b.String()
	assert.Equal(t, strings.Count(logs, fmt.Sprintf("Listener `0.0.0.0:%d` shutting down", port)), 1, logs)
	assert.Equal(t, strings.Count(logs, fmt.Sprintf("Listener `0.0.0.0:%d` stopped", port)), 1, logs)
}
