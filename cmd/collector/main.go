package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flags := ParseFlags()

	config, err := LoadConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalln("Error loading config:", err)
		return
	}

	app := &application{
		flags:  flags,
		config: config,
	}

	ctx, cancel := context.WithCancel(context.Background())
	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-osCh
		cancel()
	}()

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application error: %v\n", err)
	}
}
