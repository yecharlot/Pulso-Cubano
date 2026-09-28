package metrics

import "testing"

func TestPulseZeroWithoutObs(t *testing.T) {
	p := ComputePulse(80, 50, 70, 60, 55, 75, 70, 0)
	if p.Overall != 0 || p.ObservationN != 0 {
		t.Fatalf("%+v", p)
	}
}

func TestPulseWeighted(t *testing.T) {
	p := ComputePulse(100, 100, 100, 100, 100, 100, 100, 10)
	if p.Overall < 99 || p.Overall > 101 {
		t.Fatalf("%v", p.Overall)
	}
	if p.Layer != "observed_market" {
		t.Fatal(p.Layer)
	}
}
