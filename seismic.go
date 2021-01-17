package goseismic

import "github.com/gorilla/websocket"

// Seismic struct is type used for receiving events from websocket
type Seismic struct {
	conn      *websocket.Conn
	connected bool
	Events    chan Event
}

// New creates new Seismic value which contains Event channel for receiving seismic events
func New() *Seismic {
	return &Seismic{
		Events: make(chan Event),
	}
}
