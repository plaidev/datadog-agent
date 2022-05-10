package netflow

import (
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewNetflowServer(t *testing.T) {
	// TODO: Add server tests
	demux := aggregator.InitTestAgentDemultiplexerWithFlushInterval(10 * time.Millisecond)
	defer demux.Stop(false)

	server, err := NewNetflowServer(demux)
	require.NoError(t, err, "cannot start Netflow Server")
	assert.NotNil(t, server)
}
