package core

import (
	"encoding/json"
	"os"
	"time"
)

// Daemon launches a daemon
func (c *Controller) Daemon(exit chan<- struct{}) {
	rawMessages := make(chan [2]string)
	toLogMessages := make(chan Message)

	queryFinished := make(chan struct{})
	readFinished := make(chan struct{})
	logFinished := make(chan struct{})

	// our producer, reads from database at regular intervals
	queryTicker := time.NewTicker(time.Duration(c.cfg.DaemonInterval) * time.Second)
	go func() {
		for {
			select {
			case <-queryTicker.C:
				c.Logger.Print("Daemon: fetching new messages from DB")
				lastTS, messagesCount, err := c.queryMessages(rawMessages)
				if err != nil {
					c.Logger.Printf("Daemon: error while fetching messages: %s", err)
				}
				if lastTS != "" { // we got at least one new message, even if query failed after
					c.Logger.Printf("Daemon: %d new messages received, last timestamp is %s", messagesCount, lastTS)
					c.updateLastTimestamp(lastTS)
				}
			case <-c.ctx.Done():
				c.Logger.Print("Daemon: stopping query goroutine")
				close(queryFinished)
				return
			}
		}
	}()

	// our intermediate consumer, read messages from query producer, treat the data and dispatch it
	go func() {
		for {
			select {
			case result := <-rawMessages:
				message, err := ProcessMessage(result[0], result[1])
				if err != nil {
					c.Logger.Printf("Cannot process message with timestamp %s: %s", result[0], err)
				}
				toLogMessages <- message
			case <-queryFinished: // we stop reading only when we're sure there won't be anymore queries sent
				c.Logger.Print("Daemon: stopping processing goroutine")
				close(readFinished)
				return
			}
		}
	}()

	// a final consumer, read messages from processing producer and log it
	go func() {
		f, _ := os.OpenFile("log.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		defer f.Close()
		for {
			select {
			case message := <-toLogMessages:
				marshaledJSON, err := json.Marshal(message)
				if err != nil {
					c.Logger.Printf("Cannot marshal message: %s", err)
				}
				f.WriteString(string(marshaledJSON) + "\n")
			case <-readFinished: // we stop logging only when we're sure there won't be anymore queries processed
				f.Sync()
				c.Logger.Println("Daemon: stopping logging goroutine")
				close(logFinished)
				return
			}
		}
	}()

	<-queryFinished
	c.Logger.Print("Daemon: query goroutine is finished.")
	<-readFinished
	c.Logger.Print("Daemon: reading goroutine is finished.")
	<-logFinished
	c.Logger.Print("Daemon: logging goroutine is finished.")
	c.Logger.Print("Daemon: exited.")
	close(exit)
}
