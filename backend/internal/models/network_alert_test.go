package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNetworkAlert_JSONRoundTrip(t *testing.T) {
	raw := `{
  "agent_id": "test-host",
  "exported_at": "2025-04-01T14:22:00Z",
  "event": {
    "timestamp": "2025-04-01T14:22:00Z",
    "protocol": "TCP",
    "src_ip": "192.168.1.10",
    "src_port": 443,
    "dst_ip": "10.0.0.5",
    "dst_port": 52431,
    "watchlist_match": "dst",
    "port_name": "unknown",
    "direction": "inbound",
    "tcp_flags": "S",
    "packet_size": 128
  }
}`

	var out NetworkAlert
	require.NoError(t, json.Unmarshal([]byte(raw), &out))

	require.Equal(t, "test-host", out.AgentID)
	require.Equal(t, time.Date(2025, 4, 1, 14, 22, 0, 0, time.UTC), out.ExportedAt.UTC())
	require.Equal(t, "TCP", out.Event.Protocol)
	require.Equal(t, "192.168.1.10", out.Event.SrcIP)
	require.Equal(t, 443, *out.Event.SrcPort)
	require.Equal(t, 52431, *out.Event.DstPort)
	require.Equal(t, 128, out.Event.PacketSize)
	require.NotNil(t, out.Event.PortName)
	require.Equal(t, "unknown", *out.Event.PortName)

	b, err := json.Marshal(&out)
	require.NoError(t, err)

	var again NetworkAlert
	require.NoError(t, json.Unmarshal(b, &again))
	require.Equal(t, out.AgentID, again.AgentID)
}

// Sample from a real agent export (no port_name field).
func TestNetworkAlert_UnmarshalWithoutPortName(t *testing.T) {
	raw := `{
  "agent_id": "DESKTOP-6HLKKS8",
  "exported_at": "2026-04-01T06:24:17Z",
  "event": {
    "timestamp": "2026-04-01T06:24:17Z",
    "protocol": "TCP",
    "src_ip": "86.121.249.10",
    "src_port": 443,
    "dst_ip": "192.168.157.128",
    "dst_port": 49787,
    "watchlist_match": "dst",
    "direction": "outbound",
    "tcp_flags": "A",
    "packet_size": 60
  }
}`

	var out NetworkAlert
	require.NoError(t, json.Unmarshal([]byte(raw), &out))
	require.Equal(t, "DESKTOP-6HLKKS8", out.AgentID)
	require.Nil(t, out.Event.PortName)
	require.Equal(t, "outbound", out.Event.Direction)
	require.Equal(t, "A", out.Event.TCPFlags)
	require.Equal(t, 60, out.Event.PacketSize)
}
