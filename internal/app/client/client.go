package client

import (
	"context"
	"fmt"
	"github.com/oitimon/fawy-server/internal/app"
	"github.com/oitimon/fawy-server/pkg/pow"
	"log"
	"net"
	"sync"
	"time"
)

var requestsBuffer chan struct{}
var requestsBufferOnce sync.Once

// Client ...
type Client struct {
	config *Config
	ctx    context.Context
	wg     sync.WaitGroup
}

func NewClient(ctx context.Context, c *Config) *Client {
	return &Client{
		ctx:    ctx,
		config: c,
	}
}

func (c *Client) MultiRequests(count int) error {
	// Init pipeline for requests.
	requestsBufferOnce.Do(func() {
		requestsBuffer = make(chan struct{}, c.config.MaxRequests)
	})
	for ; count > 1; count-- {
		requestsBuffer <- struct{}{}

		select {
		case <-c.ctx.Done():
			<-requestsBuffer
			return nil
		default:
			// Create request context with timeout from Config.
			c.wg.Add(1)
			go func() {
				ctxRequest, cancelRequest := context.WithTimeout(c.ctx, time.Duration(c.config.Timeout)*time.Second)
				if err := c.Request(ctxRequest); err != nil {
					if ctxRequest.Err() == context.DeadlineExceeded {
						log.Println("Request timeout exceeded")
					} else {
						log.Printf("Client error: %v\n", err)
					}
				}
				cancelRequest()
				c.wg.Done()
			}()
		}

		<-requestsBuffer
	}

	c.wg.Wait()
	return nil
}

func (c *Client) Request(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	var conn net.Conn
	go func() {
		<-ctx.Done()
		if conn != nil {
			_ = conn.Close()
		}
		// Reset error when request was closed
		err = nil
	}()

	// Connect with server
	conn, err = net.Dial("tcp", c.config.Host+":"+c.config.Port)
	if err != nil {
		err = fmt.Errorf("connection to server error: %w", err)
		return
	}

	// Send command
	if _, err = conn.Write([]byte(app.WowCommandQuote)); err != nil {
		err = fmt.Errorf("sending error: %w", err)
		return
	}

	// Get challenge
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		err = fmt.Errorf("reading error: %w", err)
		return
	}
	challengeRequest := buffer[:n]
	log.Printf("Got challenge: %s\n", challengeRequest)

	var challenge pow.Challenge
	if challenge, err = pow.NewChallenge(c.config.Challenge); err != nil {
		return
	}
	now := time.Now()
	output, err := challenge.Fulfil(challengeRequest)
	if err != nil {
		return
	}
	log.Printf("Working time: %v", time.Since(now))
	if _, err = conn.Write(output); err != nil {
		err = fmt.Errorf("sending error: %w", err)
		return
	}

	// Read content
	if n, err = conn.Read(buffer); err != nil {
		err = fmt.Errorf("reading error: %w", err)
		return
	}
	log.Printf("Quote: %s\n\n", buffer[:n])

	return
}
