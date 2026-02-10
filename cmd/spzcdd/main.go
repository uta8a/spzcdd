package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uta8a/spzcdd/internal/app"
)

func main() {
	var listenAddr string
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:8080", "listen address")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := app.Run(runCtx, app.Config{
		ListenAddr:      listenAddr,
		ShutdownTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}
