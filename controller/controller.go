package controller

import (
	"log"
	"synrex/config"
)

// Controller is the main synrex controller
type Controller struct {
	cfg    config.Config
	logger *log.Logger
}

// New returns a new controller instance
func New(cfg config.Config, logger *log.Logger) (c Controller, err error) {
	c = Controller{
		cfg:    cfg,
		logger: logger,
	}
	return
}
