package quotes

import (
	"bytes"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository()

	if _, err := repo.Get(); err == nil || err.Error() != "repository is not filled yet" {
		t.Error("Can get values before filling")
	}

	if err := repo.Fill([][]byte{
		[]byte("Quote1"),
		[]byte("Quote2"),
		[]byte("Quote3"),
	}); err != nil {
		t.Error(err)
	}

	if v, err := repo.Get(); err != nil {
		t.Error(err)
	} else if !bytes.Equal(v, []byte("Quote1")) && !bytes.Equal(v, []byte("Quote2")) && !bytes.Equal(v, []byte("Quote3")) {
		t.Errorf("Unexpected value: %s", string(v))
	}

	// Concurrency check and random values.
	var q1, q2, q3, qN int32
	counter := 100000 // 100K
	wg := sync.WaitGroup{}
	for ; counter > 0; counter-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := repo.Get()
			if err != nil {
				t.Error(err)
			}
			if bytes.Equal(v, []byte("Quote1")) {
				atomic.AddInt32(&q1, 1)
			} else if bytes.Equal(v, []byte("Quote2")) {
				atomic.AddInt32(&q2, 1)
			} else if bytes.Equal(v, []byte("Quote3")) {
				atomic.AddInt32(&q3, 1)
			} else {
				atomic.AddInt32(&qN, 1)
			}
		}()
	}
	wg.Wait()
	if q1 == 0 || q2 == 0 || q3 == 0 || qN > 0 {
		t.Errorf("Q1: %d, Q2: %d, Q3: %d, QN: %d", q1, q2, q3, qN)
	}
}
