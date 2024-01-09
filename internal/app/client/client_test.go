package client

import (
	"context"
	"fmt"
	"github.com/oitimon/fawy-server/internal/app/quotes"
	"github.com/oitimon/fawy-server/internal/app/server"
	"github.com/oitimon/fawy-server/pkg/metrics"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	// Start server first
	ctx := context.Background()
	repo := quotes.NewRepository()
	if err := repo.Fill(quotes.GenerateQuotes()); err != nil {
		t.Error(err)
	}

	go func() {
		s, err := server.NewServer(ctx, &server.Config{
			Network:     "tcp4",
			Host:        "localhost",
			Port:        "8888",
			Timeout:     10,
			MaxHandlers: 1000,
			Difficulty:  35,
			Challenge:   "HASHBASED",
		}, metrics.NewMetrics(ctx), repo)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(s.Run())
	}()

	// Prepare Client and run single request.
	c := &Config{
		Host:        "localhost",
		Port:        "8888",
		Timeout:     5,
		MaxRequests: 100,
		Challenge:   "HASHBASED",
	}
	client := NewClient(ctx, c)
	ctxRequest, cancelRequest := context.WithTimeout(ctx, time.Duration(c.Timeout)*time.Second)
	if err := client.Request(ctxRequest); err != nil {
		if ctxRequest.Err() == context.DeadlineExceeded {
			t.Error("Request timeout exceeded")
		} else {
			t.Error("Client error:", err)
		}
	}
	cancelRequest()

	// Run 100 requests.
	if err := client.MultiRequests(100); err != nil {
		t.Error(err)
	}
}
