package goseismic

import (
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const seismicURLHost = "www.seismicportal.eu"
const seismicURLPath = "/standing_order/websocket"
const reconnectWait = 60 * time.Second

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

// Connect connects to Seismic portal websocket
func (s *Seismic) Connect() {
	u := url.URL{Scheme: "wss", Host: seismicURLHost, Path: seismicURLPath}

	var err error
	for {
		s.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err == nil {
			s.connected = true
			return
		}
		time.Sleep(reconnectWait)
	}
}

// ReadMessages reads new events (json) from seismic portal websocket, parse it and sends to channel
func (s *Seismic) ReadMessages() {
	go func() {
		for {
			if !s.connected {
				s.Connect()
			}
			_, message, err := s.conn.ReadMessage()
			if err == nil {
				if event, err := ParseEvent(message); err == nil {
					s.Events <- event
				}
			} else {
				s.Disconnect()
			}
		}
	}()
}

// Disconnect disconnects from Seismic portal websocket
func (s *Seismic) Disconnect() error {
	s.connected = false
	return s.conn.Close()
}
