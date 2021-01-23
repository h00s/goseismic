package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/h00s/goseismic"
)

func main() {
	s := goseismic.NewSeismic()

	// Set Debug field to true if you want to log what's happening behind the scenes:
	// connecting / disconnecting from websocket, sending ping and receiving pong etc.
	// default is off (s.Debug = false)
	s.Debug = true

	// Set KeepAlive to true to keep connection alive by sending pings every minute.
	// It will also disconnect and reconnect if some error ocurred while reading message.
	// default is on (s.KeepAlive = true)
	s.KeepAlive = true

	// Connect connects to websocket and starts receiving messages when successfully connected.
	// If KeepAlive is set to true and it fails to connect, it will try to automatically reconnect
	// every 15sec, otherwise it will return error.
	s.Connect()
	defer s.Disconnect()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case e := <-s.Events:
			// process received Event here
			log.Println(e)
		case <-interrupt:
			return
		}
	}
}
