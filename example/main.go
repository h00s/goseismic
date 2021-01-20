package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/h00s/goseismic"
)

func main() {
	s := goseismic.NewSeismic()
	s.Debug = true
	s.KeepAlive = true
	s.Connect()
	defer s.Disconnect()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case e := <-s.Events:
			log.Println(e)
		case <-interrupt:
			return
		}
	}
}
