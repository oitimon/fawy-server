package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Screener is just simple metrics engine to show results on screen every second.
// For local and testing using ONLY, not optimised!
type Screener struct {
	// We can not use sync.Map as we'd like to increase values (to only store).
	sync.RWMutex
	counters map[string]int64
	gauges   map[string]int64
}

func NewScreener(ctx context.Context) *Screener {
	s := &Screener{
		counters: map[string]int64{},
		gauges:   map[string]int64{},
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Do something to save/finish.
				return
			case <-time.After(time.Second):
				// Just dump everything to std.
				fmt.Printf("Counters: %v\n", s.counters)
				fmt.Printf("Gaugas: %v\n", s.gauges)
			}
		}
	}()
	return s
}

func (s *Screener) RegisterCounter(name string) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.gauges[name]; ok {
		return fmt.Errorf("you can not register %s counter as gauge already exists", name)
	}
	if _, ok := s.counters[name]; !ok {
		s.counters[name] = 0
	}
	return nil
}

func (s *Screener) RegisterGauge(name string) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.counters[name]; ok {
		return fmt.Errorf("you can not register %s gauge as counter already exists", name)
	}
	if _, ok := s.gauges[name]; !ok {
		s.gauges[name] = 0
	}
	return nil
}

func (s *Screener) Set(name string, value int64) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.counters[name]; ok {
		s.counters[name] = value
	} else if _, ok := s.gauges[name]; ok {
		s.gauges[name] = value
	}
}

func (s *Screener) Add(name string, value int64) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.counters[name]; ok {
		s.counters[name] += value
	} else if _, ok := s.gauges[name]; ok {
		s.gauges[name] += value
	}
}

func (s *Screener) Inc(name string) {
	s.Add(name, 1)
}

func (s *Screener) Get(name string) int64 {
	s.RLock()
	defer s.RUnlock()
	if v, ok := s.counters[name]; ok {
		return v
	}
	if v, ok := s.gauges[name]; ok {
		return v
	}
	return 0
}
