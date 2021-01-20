package goseismic

import "testing"

func TestNew(t *testing.T) {
	s := NewSeismic()
	if s == nil {
		t.Errorf("Seismic is nil on New()")
	}
	s.KeepAlive = false
	s.Connect()
	s.Disconnect()
}
