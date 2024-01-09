package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/oitimon/fawy-server/internal/app/client"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration.
	if _, err := os.Stat(".env_client"); err == nil {
		if err = godotenv.Load(".env_client"); err != nil {
			log.Fatal(err)
		}
	}
	var c client.Config
	if err := envconfig.Process("WOWS", &c); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Prepare Client and run single request.
	r := client.NewClient(ctx, &c)
	ctxRequest, cancelRequest := context.WithTimeout(ctx, time.Duration(c.Timeout)*time.Second)
	if err := r.Request(ctxRequest); err != nil {
		if ctxRequest.Err() == context.DeadlineExceeded {
			log.Println("Request timeout exceeded")
		} else {
			log.Printf("Client error: %v\n", err)
		}
	}
	cancelRequest()

	// Run 1000 requests.
	if err := r.MultiRequests(1000); err != nil {
		log.Println(err)
	}

	cancel()

	// Wait signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v\n", sig)
		cancel()
	case <-ctx.Done():
		// Was canceled or finished, we don't have to wait it.
		// Can be some jon here (sending buffered logs, etc).
	}
}
