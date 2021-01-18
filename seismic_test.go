package goseismic

import "testing"

func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Errorf("Seismic is nil on New()")
	}
	s.Connect()
	s.Disconnect()
}
