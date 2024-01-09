package pow

import "github.com/oitimon/fawy-server/pkg/challenge"

type Challenge interface {
	challenge.Challenge
	Difficulty
}
