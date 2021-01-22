# goseismic

[![Travis CI Build Status](https://travis-ci.com/h00s/goseismic.svg?branch=master)](https://travis-ci.com/github/h00s/goseismic)

goseismic is library for receiving (near) realtime notifications about earthquakes using websockets from [Seismicportal](https://www.seismicportal.eu/realtime.html).
JSON message is received, parsed to `goseismic.Event` and sent to channel when an event is inserted or updated. Depending on the event, you can use bots, push notification etc. to
further process the information.

## Installation

go get -u github.com/h00s/goseismic

## Usage

Information about earthquake events are sent in JSON thru websockets. This is example of one received event:

```json
{
  "action":"update",
  "data":{
    "geometry":{
      "type":"Point",
      "coordinates":[
        -121.2,
        36.6,
        -4.0
      ]
    },
    "type":"Feature",
    "id":"20201230_0000082",
    "properties":{
      "lastupdate":"2020-12-30T08:47:00.0Z",
      "magtype":"md",
      "evtype":"ke",
      "lon":-121.2,
      "auth":"NC",
      "lat":36.6,
      "depth":4.0,
      "unid":"20201230_0000082",
      "mag":2.4,
      "time":"2020-12-30T08:45:29.9Z",
      "source_id":"934165",
      "source_catalog":"EMSC-RTS",
      "flynn_region":"CENTRAL CALIFORNIA"
    }
  }
}
```

Received events are parsed and sent to `Seismic.Events` channel which you can read and further process. This is simple example how to receive events and display them (check `example/main.go` for comments):

```go
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
```