package goseismic

import (
	"log"
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
	Debug     bool
	Events    chan Event
}

// NewSeismic creates new Seismic value which contains Event channel for receiving seismic events
func NewSeismic() *Seismic {
	s := &Seismic{
		KeepAlive: true,
		Debug:     false,
		Events:    make(chan Event),
	}
	go s.sendPings()
	return s
}

// Connect connects to Seismic portal websocket
func (s *Seismic) Connect() error {
	u := url.URL{Scheme: "wss", Host: seismicURLHost, Path: seismicURLPath}

	var err error
	for {
		s.log("Connecting to websocket...")
		s.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err == nil {
			s.conn.SetPongHandler(s.pongHandler)
			s.conn.SetCloseHandler(s.closeHandler)
			s.connected = true
			go s.readMessages()
			s.log("Connected to websocket.")
			return nil
		}
		if !s.KeepAlive {
			break
		}
		time.Sleep(reconnectWait)
	}
	s.log("Error while connecting to websocket:", err)
	return err
}

// ReadMessages reads new events (JSON) from seismic portal websocket, parse it and sends to channel
func (s *Seismic) readMessages() {
	for {
		if !s.connected {
			break
		}
		_, message, err := s.conn.ReadMessage()
		if err == nil {
			if event, err := ParseEvent(message); err == nil {
				s.Events <- event
			} else {
				s.log("Error while parsing JSON")
			}
		} else {
			s.log("Error while reading message from websocket:", err)
			s.Disconnect()
			break
		}
	}
}

// Disconnect disconnects from Seismic portal websocket
func (s *Seismic) Disconnect() error {
	if s.connected {
		s.connected = false
		s.log("Disconnecting from websocket...")
		return s.conn.Close()
	}
	return nil
}

// sendPings sends control messages (ping) every pingWait interval to keep connection alive
func (s *Seismic) sendPings() {
	ticker := time.NewTicker(pingWait)
	defer ticker.Stop()

	for {
		<-ticker.C
		if s.KeepAlive {
			if s.connected {
				s.log("Sending ping...")
				s.conn.WriteMessage(websocket.PingMessage, []byte{})
			} else {
				s.Connect()
			}
		}
	}
}

func (s *Seismic) pongHandler(string) error {
	s.log("Received pong.")
	return nil
}

func (s *Seismic) closeHandler(code int, text string) error {
	s.log("Close handler called.")
	s.connected = false
	return nil
}

func (s *Seismic) log(messages ...interface{}) {
	if s.Debug {
		log.Println(messages...)
	}
}
