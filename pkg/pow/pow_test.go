package pow

import (
	"reflect"
	"testing"
)

func TestNewGoPow(t *testing.T) {
	ch, err := NewChallenge("GO-POW")
	if err != nil {
		t.Error("Can not create GO-POW challenge")
	}
	if reflect.TypeOf(ch) != reflect.TypeOf(&GoPow{}) {
		t.Errorf("Expected %s type of challenge but got %s", reflect.TypeOf(&GoPow{}), reflect.TypeOf(ch))
	}
}

func TestNewNumeric(t *testing.T) {
	ch, err := NewChallenge("NUMERIC")
	if err != nil {
		t.Error("Can not create NUMERIC challenge")
	}
	if reflect.TypeOf(ch) != reflect.TypeOf(&Numeric{}) {
		t.Errorf("Expected %s type of challenge but got %s", reflect.TypeOf(&Numeric{}), reflect.TypeOf(ch))
	}
}

func TestNewHASHBASED(t *testing.T) {
	ch, err := NewChallenge("HASHBASED")
	if err != nil {
		t.Error("Can not create HASHBASED challenge")
	}
	if reflect.TypeOf(ch) != reflect.TypeOf(&Hashbased{}) {
		t.Errorf("Expected %s type of challenge but got %s", reflect.TypeOf(&Hashbased{}), reflect.TypeOf(ch))
	}
}

func TestNewUndefined(t *testing.T) {
	_, err := NewChallenge("something")
	if err == nil {
		t.Error("Error expected but challenge was created")
	}
}
