package goseismic

import (
	"encoding/json"
	"log"
	"time"
)

// Event represent one information obtained from seismic portal
/*{
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
*/
type Event struct {
	Action string `json:"action"`
	Data   struct {
		Geometry struct {
			Type        string     `json:"type"`
			Coordinates [3]float64 `json:"coordinates"`
		} `json:"geometry"`
		Type       string `json:"type"`
		ID         string `json:"id"`
		Properties struct {
			LastUpdate    time.Time `json:"lastupdate"`
			MagType       string    `json:"magtype"`
			EvType        string    `json:"evtype"`
			Longitude     float64   `json:"lon"`
			Auth          string    `json:"auth"`
			Latitude      float64   `json:"lat"`
			Depth         float64   `json:"depth"`
			UnID          string    `json:"unid"`
			Magnitude     float64   `json:"mag"`
			Time          time.Time `json:"time"`
			SourceID      string    `json:"source_id"`
			SourceCatalog string    `json:"source_catalog"`
			FlynnRegion   string    `json:"flynn_region"`
		} `json:"properties"`
	} `json:"data"`
}

// ParseEvent converts received message from websocket to Event struct
func ParseEvent(data []byte) (Event, error) {
	var event Event

	if err := json.Unmarshal(data, &event); err != nil {
		log.Println("Error parsing JSON")
		return event, err
	}

	return event, nil
}
