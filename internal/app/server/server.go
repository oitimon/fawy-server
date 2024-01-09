package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/oitimon/fawy-server/internal/app"
	"github.com/oitimon/fawy-server/internal/app/quotes"
	"github.com/oitimon/fawy-server/pkg/metrics"
	"github.com/oitimon/fawy-server/pkg/pow"
	"github.com/valyala/tcplisten"
	"log"
	"net"
	"sync"
	"time"
)

var handlersBuffer chan struct{}
var handlersBufferOnce sync.Once

// Server ...
type Server struct {
	config  *Config
	metrics metrics.Metrics
	repo    quotes.Repository
	ctx     context.Context
	wg      sync.WaitGroup
}

func NewServer(ctx context.Context, c *Config, m metrics.Metrics, r quotes.Repository) (s *Server, err error) {
	if err = m.RegisterCounter(MetricProceed); err != nil {
		return
	}
	if err = m.RegisterCounter(MetricErrors); err != nil {
		return
	}
	if err = m.RegisterCounter(MetricTimeouts); err != nil {
		return
	}
	if err = m.RegisterGauge(MetricHandling); err != nil {
		return
	}
	s = &Server{
		config:  c,
		metrics: m,
		repo:    r,
		ctx:     ctx,
	}
	return
}

// Run ...
func (s *Server) Run() error {
	log.Printf("Server started as %s:%s", s.config.Host, s.config.Port)

	tcpConfig := &tcplisten.Config{
		ReusePort: true,
		FastOpen:  true,
	}
	listener, err := tcpConfig.NewListener(s.config.Network, fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))
	if err != nil {
		return err
	}
	defer func(listener net.Listener) {
		if err := listener.Close(); err != nil {
			log.Println(err)
		}
	}(listener)

	// Init pipeline for handlers.
	handlersBufferOnce.Do(func() {
		handlersBuffer = make(chan struct{}, s.config.MaxHandlers)
	})

	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				return err
			}
			s.wg.Add(1)
			// Create handler's context with timeout from Config.
			ctxHandler, cancel := context.WithTimeout(s.ctx, time.Duration(s.config.Timeout)*time.Second)
			go func() {
				err = s.handleRequest(ctxHandler, conn)
				if ctxHandler.Err() == context.DeadlineExceeded {
					s.metrics.Inc(MetricTimeouts)
					//@todo debug
					log.Println("Handler timeout")
				} else if ctxHandler.Err() == nil && err != nil {
					s.metrics.Inc(MetricErrors)
					//@todo debug
					//log.Println("Error:", err)
				}
				cancel()
				s.wg.Done()
			}()
		}
	}

}

func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) handleRequest(ctx context.Context, conn net.Conn) error {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
	}()
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Wait for pipeline capacity.
	handlersBuffer <- struct{}{}
	defer func() {
		<-handlersBuffer
	}()

	s.metrics.Inc(MetricHandling)
	defer func() {
		s.metrics.Add(MetricHandling, -1)
		s.metrics.Inc(MetricProceed)
	}()

	// Get and validate COMMAND request
	buffer := make([]byte, 32)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("reading error: %w", err)
	}
	if !bytes.Equal(buffer[:n], []byte(app.WowCommandQuote)) {
		if _, err = conn.Write([]byte("Not correct command")); err != nil {
			return fmt.Errorf("sending wrong command message error: %w", err)
		}
		return fmt.Errorf("not correct command from client - %s", buffer[:n])
	}

	// Create challenge
	var challenge pow.Challenge
	if challenge, err = pow.NewChallenge(s.config.Challenge); err != nil {
		return err
	}
	challenge.SetDifficulty(s.config.Difficulty)
	challengeRequest, err := challenge.Request()
	if err != nil {
		return err
	}

	_, err = conn.Write(challengeRequest)
	if err != nil {
		return fmt.Errorf("sending error: %w", err)
	}

	// Get challenge response
	if n, err = conn.Read(buffer); err != nil {
		return fmt.Errorf("reading error: %w", err)
	}

	// Check challenge
	ok, err := challenge.Check(buffer[:n])
	if err != nil {
		return fmt.Errorf("challenge checking error: %w", err)
	}
	if !ok {
		return fmt.Errorf("challenge is not accepted: %s", string(buffer[:n]))
	}

	// Generate content @todo quotes repository
	//content := "I'm the super-duper quote"
	content, err := s.repo.Get()
	if err != nil {
		return fmt.Errorf("error getting quote from repository: %w", err)
	}
	_, err = conn.Write(content)
	if err != nil {
		return fmt.Errorf("sending error: %w", err)
	}

	return nil
}
