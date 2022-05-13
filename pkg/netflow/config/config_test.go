package config

import (
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/netflow/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestReadConfig(t *testing.T) {
	var tests = []struct {
		name           string
		configYaml     string
		expectedConfig NetflowConfig
		expectedError  string
	}{
		{
			name: "basic configs",
			configYaml: `
network_devices:
  netflow:
    enabled: true
    stop_timeout: 10
    aggregator_buffer_size: 20
    aggregator_flush_interval: 30
    log_payloads: true
    listeners:
      - flow_type: netflow9
        bind_host: 127.0.0.1
        port: 1234
`,
			expectedConfig: NetflowConfig{
				StopTimeout:             10,
				AggregatorBufferSize:    20,
				AggregatorFlushInterval: 30,
				LogPayloads:             true,
				Listeners: []ListenerConfig{
					{
						FlowType: common.TypeNetFlow9,
						BindHost: "127.0.0.1",
						Port:     uint16(1234),
					},
				},
			},
		},
		{
			name: "defaults",
			configYaml: `
network_devices:
  netflow:
    enabled: true
    listeners:
      - flow_type: netflow9
`,
			expectedConfig: NetflowConfig{
				StopTimeout:             5,
				AggregatorBufferSize:    100,
				AggregatorFlushInterval: 10,
				LogPayloads:             false,
				Listeners: []ListenerConfig{
					{
						FlowType: common.TypeNetFlow9,
						BindHost: "localhost",
						Port:     uint16(2055),
					},
				},
			},
		},
		{
			name: "invalid flow type",
			configYaml: `
network_devices:
  netflow:
    enabled: true
    listeners:
      - flow_type: invalidType
`,
			expectedError: "the provided flow type `invalidType` is not valid (valid flow types: [ipfix sflow5 netflow5 netflow9])",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Datadog.SetConfigType("yaml")
			err := config.Datadog.ReadConfig(strings.NewReader(tt.configYaml))
			require.NoError(t, err)

			readConfig, err := ReadConfig()
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, readConfig)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedConfig, *readConfig)
			}
		})
	}
}

func TestListenerConfig_Addr(t *testing.T) {
	listenerConfig := ListenerConfig{
		FlowType: common.TypeNetFlow9,
		BindHost: "127.0.0.1",
		Port:     1234,
	}
	assert.Equal(t, "127.0.0.1:1234", listenerConfig.Addr())
}
