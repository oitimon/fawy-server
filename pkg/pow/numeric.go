package pow

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Numeric - Challenge simple implementation (no real POW, just example).
type Numeric struct {
	difficulty uint
	request    int
}

func NewNumeric() *Numeric {
	p := &Numeric{}
	return p
}

func (p *Numeric) SetDifficulty(d uint) {
	p.difficulty = d
	if p.difficulty > 100 {
		p.difficulty = 100
	}
}

func (p *Numeric) Request() ([]byte, error) {
	p.request = rand.Intn(100000)
	return []byte(strconv.Itoa(p.request)), nil
}

func (p *Numeric) Fulfil(request []byte) (output []byte, err error) {
	i, err := strconv.Atoi(string(request))
	if err != nil {
		err = fmt.Errorf("challenge request in not INT: %s", string(request))
		return
	}
	// Super MATH.
	output = []byte(strconv.Itoa(i + 2))
	// Emulate "working" by pausing of difficulty.
	time.Sleep(time.Duration(p.difficulty) * time.Millisecond)
	return
}

func (p *Numeric) Check(result []byte) (ok bool, err error) {
	i, err := strconv.Atoi(string(result))
	if err != nil {
		err = fmt.Errorf("challenge result in not INT: %s", string(result))
		return
	}
	// Validate Super MATH.
	ok = i-2 == p.request
	return
}
