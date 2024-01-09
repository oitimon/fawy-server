package pow

import (
	"fmt"
)

func NewChallenge(name string) (Challenge, error) {
	switch name {
	case "GO-POW":
		return NewGoPow(), nil
	case "NUMERIC":
		return NewNumeric(), nil
	case "HASHBASED":
		return NewHashbased(), nil
	default:
		return nil, fmt.Errorf("undefined pow name %s", name)
	}
}
