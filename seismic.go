package goseismic

import "github.com/gorilla/websocket"

// Seismic struct is type used for receiving events from websocket
type Seismic struct {
	conn      *websocket.Conn
	connected bool
	Events    chan Event
}
