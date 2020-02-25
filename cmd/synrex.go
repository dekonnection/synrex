package main

import (
	"fmt"
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

	c, err := controller.New(cfg)

	fmt.Println(c)
}
