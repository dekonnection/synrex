package core

import (
	"context"
	"log"
	"synrex/config"
)

// Controller is the main synrex controller
type Controller struct {
	ctx           context.Context
	cfg           config.Config
	Logger        *log.Logger
	roomsList     []string
	roomsIndex    map[string]string
	rawMessages   chan [2]string
	lastTimestamp string
}

// NewController returns a new controller instance
func NewController(ctx context.Context, logger *log.Logger, cfg config.Config) (c Controller, err error) {
	c = Controller{
		ctx:        ctx,
		cfg:        cfg,
		Logger:     logger,
		roomsList:  []string{},
		roomsIndex: map[string]string{},
		// rawMessages: make(chan [2]string),
	}
	for roomName, roomID := range cfg.Rooms {
		c.roomsList = append(c.roomsList, roomID)
		c.roomsIndex[roomID] = roomName
	}

	err = c.readLastTimestamp()
	if err != nil {
		c.Logger.Print("Cannot read last timestamp.")
		err = c.updateLastTimestamp("0")
		if err != nil {
			return
		}
	}
	c.Logger.Printf("Last timestamp is: %s", c.lastTimestamp)
	c.Logger.Print("Controller initialized.")

	return
}
