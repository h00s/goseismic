package goseismic

import "testing"

func TestNew(t *testing.T) {
	if s := New(); s == nil {
		t.Errorf("Seismic is nil on New()")
	}
}
