package netflow

import (
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/aggregator/mocksender"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestNewNetflowServer(t *testing.T) {
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

	// Send netflowV5Data twice to test aggregator
	// Flows will have 2x bytes/packets after aggregation
	err = sendUDPPacket(port, mockNetflowV5Data)
	require.NoError(t, err)
	err = sendUDPPacket(port, mockNetflowV5Data)
	require.NoError(t, err)

	// Make sure flows are flushed before doing assertions
	waitFlowsToBeFlushed(server.flowAgg, 10*time.Second)

	sender.AssertEventPlatformEvent(t, mock.Anything, "network-devices-netflow")
	sender.AssertMetric(t, "Count", "datadog.newflow.aggregator.flows_received", 1, "", []string{"sample_addr:127.0.0.1", "flow_type:netflow5"})
}

func TestStartServerAndStopServer(t *testing.T) {
	demux := aggregator.InitTestAgentDemultiplexerWithFlushInterval(10 * time.Millisecond)
	defer demux.Stop(false)
	err := StartServer(demux)
	require.NoError(t, err)
	require.NotNil(t, serverInstance)

	replaceWithDummyFlowProcessor(serverInstance, 123)

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
	// Setup NetFlow config
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

	// Setup Netflow Server
	demux := aggregator.InitTestAgentDemultiplexerWithFlushInterval(10 * time.Millisecond)
	defer demux.Stop(false)
	server, err := NewNetflowServer(demux)
	require.NoError(t, err, "cannot start Netflow Server")
	assert.NotNil(t, server)

	flowProcessor := replaceWithDummyFlowProcessor(server, port)

	// Stops server
	server.stop()

	// Assert logs present
	assert.Equal(t, flowProcessor.stopped, true)
}
