package pow

type Difficulty interface {
	// SetDifficulty must be 1-100.
	SetDifficulty(uint)
}
