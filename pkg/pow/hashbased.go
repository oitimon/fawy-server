package pow

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"
	"time"
)

const hexCharset = "abcdef0123456789"

type Hashbased struct {
	difficulty   uint
	targetPrefix []byte
	data         []byte
}

func NewHashbased() *Hashbased {
	p := &Hashbased{}
	return p
}

func (p *Hashbased) SetDifficulty(d uint) {
	// For Hashbased, we have 12 as maximum difficulty.
	p.difficulty = d / 8
	if p.difficulty > 12 {
		p.difficulty = 12
	} else if p.difficulty < 1 {
		p.difficulty = 1
	}
}

func (p *Hashbased) Request() ([]byte, error) {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	p.targetPrefix = make([]byte, p.difficulty)
	for i := range p.targetPrefix {
		p.targetPrefix[i] = hexCharset[seededRand.Intn(16)]
	}

	p.data = make([]byte, 12)
	for i := range p.data {
		p.data[i] = hexCharset[seededRand.Intn(16)]
	}

	// As request, we have prefix in necessary length (difficulty) and data colon separated.
	return append(append(p.targetPrefix, ':'), p.data...), nil
}

func (p *Hashbased) Fulfil(request []byte) (output []byte, err error) {
	// Split request for target prefix and data
	rr := bytes.Split(request, []byte(":"))
	if len(rr) != 2 {
		err = errors.New("hashbased request must have two parts splitted by colon")
		return
	}
	targetPrefix := rr[0]
	data := rr[1]
	difficulty := len(targetPrefix)

	// Find the solution.
	var nonce int64
	for {
		attempt := append(data, []byte(strconv.FormatInt(nonce, 10))...)
		hash := p.calculateHash(attempt)

		if bytes.Equal(hash[:difficulty], targetPrefix) {
			return []byte(strconv.FormatInt(nonce, 10)), nil
		}

		nonce++
	}
}

func (p *Hashbased) Check(result []byte) (ok bool, err error) {
	hash := p.calculateHash(append(p.data, result...))

	return bytes.Equal(hash[:p.difficulty], p.targetPrefix), nil
}

func (p *Hashbased) calculateHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	hashEncoded := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(hashEncoded, hash[:])
	return hashEncoded
}
