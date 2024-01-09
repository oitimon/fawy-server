package pow

import (
	"testing"
	"time"
)

func TestHashbased(t *testing.T) {
	ch, err := NewChallenge("HASHBASED")
	if err != nil {
		t.Error("Can not create HASHBASED challenge")
	}
	ch.SetDifficulty(35)

	request, err := ch.Request()
	if err != nil {
		t.Errorf("Error when request: %v", err)
	}
	if len(request) == 0 {
		t.Error("Empty request")
	}

	// Fulfill and check time.
	n := time.Now()
	response, err := ch.Fulfil(request)
	if err != nil {
		t.Errorf("Error when fulfill: %v", err)
	}
	if len(response) == 0 {
		t.Error("Empty response")
	}
	d := time.Since(n)
	if d < time.Millisecond {
		t.Errorf("Fulfill is to fast: %v", d)
	}

	// Valid check.
	ok, err := ch.Check(response)
	if err != nil {
		t.Errorf("Error when check: %v", err)
	}
	if !ok {
		t.Error("Check is not valid when it must")
	}

	// Invalid check.
	ok, err = ch.Check([]byte("123"))
	if err != nil {
		t.Errorf("Error when check: %v", err)
	}
	if ok {
		t.Error("Check is valid when it must not")
	}
}
