package quotes

import (
	"errors"
	"math/rand"
)

type Repository interface {
	// Fill storage before using.
	Fill([][]byte) error
	// Get random value from storage.
	Get() ([]byte, error)
}

type Memory struct {
	data [][]byte
	len  int
}

// NewRepository We have only one implementation for now.
func NewRepository() Repository {
	return NewMemory()
}

func NewMemory() *Memory {
	return &Memory{
		data: [][]byte{},
	}
}

// Fill Concurrency NOT SAFE!
// Please use once after creating and before reading.
func (m *Memory) Fill(data [][]byte) error {
	m.data = data
	m.len = len(m.data)
	return nil
}

// Get Concurrency safe after filling.
func (m *Memory) Get() (value []byte, err error) {
	if m.len == 0 {
		err = errors.New("repository is not filled yet")
	} else {
		value = m.data[rand.Intn(m.len)]
	}
	return
}
