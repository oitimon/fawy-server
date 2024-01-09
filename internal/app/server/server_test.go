package server

import (
	"context"
	"fmt"
	"github.com/oitimon/fawy-server/internal/app/quotes"
	"github.com/oitimon/fawy-server/pkg/metrics"
	"github.com/oitimon/fawy-server/pkg/pow"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var once sync.Once
var benchCount, requestCount, successCount int32
var semRequestor chan struct{}

func BenchmarkTCPServer(b *testing.B) {
	once.Do(func() {
		c := &Config{
			Network:     "tcp4",
			Host:        "localhost",
			Port:        "8888",
			Timeout:     3,
			MaxHandlers: 1000,
			Difficulty:  35,
			Challenge:   "HASHBASED",
		}

		atomic.AddInt32(&benchCount, 1)
		ctx := context.Background()

		repo := quotes.NewRepository()
		if err := repo.Fill(quotes.GenerateQuotes()); err != nil {
			b.Error(err)
		}

		go func() {
			s, err := NewServer(ctx, c, metrics.NewMetrics(ctx), repo)
			if err != nil {
				b.Error(err)
			}
			fmt.Println(s.Run())
		}()

		n := 2000
		semRequestor = make(chan struct{}, 100)

		wg := sync.WaitGroup{}
		tstamp := time.Now()
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				sendRequest(b, c)
				wg.Done()
			}()
		}
		wg.Wait()

		fmt.Println("Benchmarks count: ", benchCount)
		fmt.Println("Total requests count: ", requestCount)
		fmt.Println("From them success are: ", successCount)
		fmt.Println("Time: ", time.Since(tstamp))

		if requestCount != successCount {
			b.Errorf("%d requests but only %d are success", requestCount, successCount)
		}
	})
}

func sendRequest(b *testing.B, c *Config) {
	semRequestor <- struct{}{}
	defer func() {
		<-semRequestor
	}()
	atomic.AddInt32(&requestCount, 1)

	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println("Error creating connection:", err)
		return
	}
	defer func(conn net.Conn) {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}(conn)

	// Send hello data
	message := "WOW-QUOTE"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Sending error:", err)
		return
	}

	// Get challenge
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		b.Errorf("Reading error: %v", err)
	}
	//fmt.Printf("Got challenge: %s\n", buffer[:n])

	var challenge pow.Challenge
	if challenge, err = pow.NewChallenge(c.Challenge); err != nil {
		b.Errorf("Challenge error: %v", err)
	}
	//now := time.Now()
	output, err := challenge.Fulfil(buffer[:n])
	if err != nil {
		b.Error(err)
	}
	//fmt.Println("Working time: ", time.Since(now))
	_, err = conn.Write(output)
	if err != nil {
		b.Errorf("Sending error: %v", err)
	}

	// Read content
	_, err = conn.Read(buffer)
	if err != nil {
		b.Errorf("Reading error: %v", err)
		return
	}
	//fmt.Printf("Quote: %s\n", buffer[:n])
	atomic.AddInt32(&successCount, 1)
}
