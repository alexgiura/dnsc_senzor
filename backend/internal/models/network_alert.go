package models

import "time"

type NetworkAlert struct {
	AgentID    string            `json:"agent_id"`
	ExportedAt time.Time         `json:"exported_at"`
	Event      NetworkAlertEvent `json:"event"`
}

type NetworkAlertEvent struct {
	Timestamp      time.Time `json:"timestamp"`
	Protocol       string    `json:"protocol"`
	SrcIP          string    `json:"src_ip"`
	SrcPort        *int      `json:"src_port,omitempty"`
	DstIP          string    `json:"dst_ip"`
	DstPort        *int      `json:"dst_port,omitempty"`
	WatchlistMatch string    `json:"watchlist_match"`
	PortName       *string   `json:"port_name,omitempty"`
	Direction      string    `json:"direction"`
	TCPFlags       string    `json:"tcp_flags"`
	PacketSize     int       `json:"packet_size"`
}
