package core

import (
	"fmt"
	"time"
)

// Daemon launches a daemon
func (c *Controller) Daemon(exit chan<- struct{}) {
	rawMessages := make(chan [2]string)
	queryFinished := make(chan struct{})
	readFinished := make(chan struct{})

	// our producer, reads from database at regular intervals
	queryTicker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-queryTicker.C:
				c.Logger.Print("Daemon: fetching new messages from DB")
				err := c.queryMessages(rawMessages)
				if err != nil {
					c.Logger.Printf("Daemon: cannot fetch messages: %s", err)
				}
			case <-c.ctx.Done():
				c.Logger.Print("Daemon: stopping query goroutine")
				close(queryFinished)
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case result := <-rawMessages:
				fmt.Println(result)
			case <-queryFinished: // we stop reading only when we're sure there won't be anymore queries sent
				c.Logger.Print("Daemon: stopping reading goroutine")
				close(readFinished)
				return
			}
		}
	}()
	<-queryFinished
	c.Logger.Print("Daemon: query goroutine is finished.")
	<-readFinished
	c.Logger.Print("Daemon: reading goroutine is finished.")
	c.Logger.Print("Daemon: exited.")
	close(exit)
}
