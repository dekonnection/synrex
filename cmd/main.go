package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"synrex/config"
	"synrex/core"
	"syscall"
)

var (
	exitMain  chan struct{}
	ctx       context.Context
	cancelCtx context.CancelFunc
)

func main() {
	fmt.Println("Starting synrex ...")

	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Printf("Cannot load config: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Configuration loaded.")

	logger := log.New(os.Stderr, "[Synrex] ", 3)
	logger.Print("Logging initialized.")

	signals := make(chan os.Signal, 1)
	exitMain = make(chan struct{})
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		logger.Print("Exit signal received.")
		exit()
	}()

	ctx, cancelCtx = context.WithCancel(context.Background())
	defer cancelCtx()

	c, err := core.NewController(ctx, logger, cfg)
	if err != nil {
		logger.Fatal("Cannot initialize controller, will exit.")
	}
	// initialization is done

	exitDaemon := make(chan struct{})
	go c.Daemon(exitDaemon)

	c.Logger.Print("Main: daemon is running.")
	<-exitMain
	c.Logger.Print("Main: exiting, waiting for daemon to stop")
	<-exitDaemon
	c.Logger.Print("Main: daemon stopped, bye.")

}

func exit() {
	cancelCtx()
	close(exitMain)
}
