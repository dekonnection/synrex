package main

import (
	"fmt"
	"log"
	"os"
	"synrex/config"
	"synrex/controller"
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

	c, err := controller.New(cfg, logger)
	if err != nil {
		logger.Fatal("Cannot initialize controller, will exit.")
	}
	// initialization is done

}
