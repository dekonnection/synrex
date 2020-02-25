package controller

import "synrex/config"

// Controller is the main synrex controller
type Controller struct {
	cfg config.Config
}

// New returns a new controller instance
func New(cfg config.Config) (c Controller, err error) {
	c = Controller{
		cfg: cfg,
	}
	return
}
