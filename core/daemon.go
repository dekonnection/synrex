package core

import "fmt"

// Daemon launches a daemon
func (c *Controller) Daemon() {
	rawMessages := make(chan [2]string)
	readFinished := make(chan struct{})

	go c.queryMessages(rawMessages)

	go func() {
		for result := range rawMessages {
			fmt.Println(result)
		}
		close(readFinished)
	}()
	<-readFinished
	c.Logger.Print("Daemon: Reading is finished.")
}
