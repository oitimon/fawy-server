package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/oitimon/fawy-server/internal/app/quotes"
	"github.com/oitimon/fawy-server/internal/app/server"
	"github.com/oitimon/fawy-server/pkg/metrics"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration.
	if _, err := os.Stat(".env_server"); err == nil {
		if err = godotenv.Load(".env_server"); err != nil {
			log.Fatal(err)
		}
	}
	var c server.Config
	if err := envconfig.Process("WOWS", &c); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Prepare Repository with quotes and Server.
	repo := quotes.NewRepository()
	if err := repo.Fill(quotes.GenerateQuotes()); err != nil {
		log.Fatal(err)
	}
	s, err := server.NewServer(ctx, &c, metrics.NewMetrics(ctx), repo)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err := s.Run()
		log.Println("Server error:", err)
		// We cancel if Server generates fatal in the process.
		cancel()
	}()

	// Wait signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v\n", sig)
		cancel()
		// Wait all server's handlers close connections and finish.
		s.Wait()
	case <-ctx.Done():
		// Was canceled by server, we don't have to wait it.
		// Can be some jon here (sending buffered logs, etc).
	}
}
