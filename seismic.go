package goseismic

import (
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const seismicURLHost = "www.seismicportal.eu"
const seismicURLPath = "/standing_order/websocket"
const reconnectWait = 15 * time.Second
const pingWait = 60 * time.Second

// Seismic struct is type used for receiving events from websocket
type Seismic struct {
	conn      *websocket.Conn
	connected bool
	KeepAlive bool
	Events    chan Event
}

// NewSeismic creates new Seismic value which contains Event channel for receiving seismic events
func NewSeismic() *Seismic {
	s := &Seismic{
		KeepAlive: true,
		Events:    make(chan Event),
	}
	go s.readMessages()
	go s.sendPings()
	return s
}

// Connect connects to Seismic portal websocket
func (s *Seismic) Connect() {
	u := url.URL{Scheme: "wss", Host: seismicURLHost, Path: seismicURLPath}

	var err error
	for {
		s.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err == nil {
			s.connected = true
			s.conn.SetPongHandler(s.pongHandler)
			s.conn.SetCloseHandler(s.closeHandler)
			return
		}
		time.Sleep(reconnectWait)
	}
}

// ReadMessages reads new events (json) from seismic portal websocket, parse it and sends to channel
func (s *Seismic) readMessages() {
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
}

// Disconnect disconnects from Seismic portal websocket
func (s *Seismic) Disconnect() error {
	s.connected = false
	return s.conn.Close()
}

// sendPings sends control messages (ping) every pingWait interval to keep connection alive
func (s *Seismic) sendPings() {
	ticker := time.NewTicker(pingWait)
	defer ticker.Stop()

	for {
		<-ticker.C
		if s.connected && s.KeepAlive {
			s.conn.WriteMessage(websocket.PingMessage, []byte{})
		}
	}
}

func (s *Seismic) pongHandler(string) error {
	return nil
}

func (s *Seismic) closeHandler(code int, text string) error {
	s.connected = false
	return nil
}
