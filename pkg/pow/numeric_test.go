package pow

import (
	"strconv"
	"testing"
)

func TestNumeric(t *testing.T) {
	ch, err := NewChallenge("NUMERIC")
	if err != nil {
		t.Error("Can not create NUMERIC challenge")
	}
	ch.SetDifficulty(5)

	request, err := ch.Request()
	if err != nil {
		t.Errorf("Error when request: %v", err)
	}
	if len(request) == 0 {
		t.Error("Empty request")
	}

	response, err := ch.Fulfil(request)
	if err != nil {
		t.Errorf("Error when fulfill: %v", err)
	}
	if len(response) == 0 {
		t.Error("Empty response")
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
	ok, err = ch.Check(append(response, []byte("1")...))
	if err != nil {
		t.Errorf("Error when check: %v", err)
	}
	if ok {
		t.Error("Check is valid when it must not")
	}

	// Self-check.
	requestI, err := strconv.Atoi(string(request))
	if err != nil {
		t.Errorf("Challenge request in not INT: %s", string(request))
	}
	responseI, err := strconv.Atoi(string(response))
	if err != nil {
		t.Errorf("Challenge response in not INT: %s", string(request))
	}
	if responseI != requestI+2 {
		t.Errorf("Expected response %d but got %d", requestI+2, responseI)
	}
}
