package example

import "time"

const (
	Port           = 3000
	WebsocketRoute = "ws"

	// some mysterious consts
	WriteWait = 10 * time.Second
	PongWait  = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod     = (PongWait * 1) / 10
	MaxMessageSize = 512
)
