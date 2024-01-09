package metrics

import (
	"context"
	"testing"
)

func TestScreener(t *testing.T) {
	s := NewScreener(context.Background())
	err := s.RegisterCounter("cou1")
	if err != nil {
		t.Errorf("Error when resiterging counter: %v", err)
	}

	err = s.RegisterGauge("cou1")
	if err == nil || err.Error() != "you can not register cou1 gauge as counter already exists" {
		t.Error("Could register gauge when counter is present")
	}

	err = s.RegisterGauge("gau1")
	if err != nil {
		t.Errorf("Error when resiterging gauge: %v", err)
	}

	err = s.RegisterCounter("gau1")
	if err == nil || err.Error() != "you can not register gau1 counter as gauge already exists" {
		t.Error("Could register counter when gauge is present")
	}

	s.Set("cou1", 123)
	s.Add("cou1", 10)
	s.Inc("cou1")
	if s.Get("cou1") != 134 {
		t.Errorf("Expected 134 from counter but got %d", s.Get("cou1"))
	}

	s.Set("gau1", 223)
	s.Add("gau1", -10)
	s.Inc("gau1")
	if s.Get("gau1") != 214 {
		t.Errorf("Expected 214 from gauge but got %d", s.Get("gau1"))
	}

	s.Set("test", 100) // Should be ignored.
	if s.Get("test") != 0 {
		t.Errorf("Expected 0 from test but got %d", s.Get("test"))
	}
}
