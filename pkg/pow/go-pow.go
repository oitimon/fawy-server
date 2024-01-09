package pow

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/bwesterb/go-pow"
)

// GoPow - Challenge driver for GO-POW library.
type GoPow struct {
	request    string
	bound      []byte
	difficulty uint
}

func NewGoPow() *GoPow {
	p := &GoPow{}
	return p
}

func (p *GoPow) SetDifficulty(d uint) {
	// For GoPow we have 50 as maximum difficulty and 10 as minimum.
	p.difficulty = d / 2
	if p.difficulty > 50 {
		p.difficulty = 50
	} else if p.difficulty < 10 {
		p.difficulty = 10
	}
}

func (p *GoPow) Request() (res []byte, err error) {
	nonce := make([]byte, 128)
	if _, err = rand.Read(nonce); err != nil {
		err = fmt.Errorf("can not create a nonce: %w", err)
		return
	}
	p.request = pow.NewRequest(uint32(p.difficulty), nonce)
	p.bound = make([]byte, 16)
	if _, err = rand.Read(p.bound); err != nil {
		err = fmt.Errorf("can not create a bound: %w", err)
		return
	}
	p.bound = []byte("some bound data")
	// As a challenge request we combine internal request with bound data with colon splitter.
	requestB := []byte(p.request)
	reqEncoded := make([]byte, base64.StdEncoding.EncodedLen(len(requestB)))
	base64.StdEncoding.Encode(reqEncoded, requestB)
	boundEncoded := make([]byte, base64.StdEncoding.EncodedLen(len(p.bound)))
	base64.StdEncoding.Encode(boundEncoded, p.bound)
	res = append(append(reqEncoded, ':'), boundEncoded...)
	return
}

func (p *GoPow) Fulfil(request []byte) (output []byte, err error) {
	// Split request for nonce-request and bound data
	rr := bytes.Split(request, []byte(":"))
	if len(rr) != 2 {
		err = errors.New("GoPow request must have two parts splitted by colon")
		return
	}
	requestPart := make([]byte, base64.StdEncoding.DecodedLen(len(rr[0])))
	if _, err = base64.StdEncoding.Decode(requestPart, rr[0]); err != nil {
		err = fmt.Errorf("can not decode request part: %w", err)
		return
	}
	bound := make([]byte, base64.StdEncoding.DecodedLen(len(rr[1])))
	if _, err = base64.StdEncoding.Decode(bound, rr[1]); err != nil {
		err = fmt.Errorf("can not decode request part: %w", err)
		return
	}
	outputS, err := pow.Fulfil(string(requestPart), bound)
	if err != nil {
		err = fmt.Errorf("can not fulfil request: %w", err)
		return
	}
	output = []byte(outputS)
	return
}

func (p *GoPow) Check(result []byte) (ok bool, err error) {
	if ok, err = pow.Check(p.request, string(result), p.bound); err != nil {
		err = fmt.Errorf("can not check response: %w", err)
	}
	return
}
