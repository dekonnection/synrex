package core

import (
	"errors"
	"io/ioutil"
)

func (c *Controller) readLastTimestamp() (err error) {
	rawTimestamp, err := ioutil.ReadFile(c.cfg.LastTsFile)
	if err != nil {
		err = errors.New("cannot open timestamp file")
		return
	}
	c.lastTimestamp = string(rawTimestamp)
	return
}

func (c *Controller) writeLastTimestamp() (err error) {
	err = ioutil.WriteFile(c.cfg.LastTsFile, []byte(c.lastTimestamp), 0600)
	if err != nil {
		err = errors.New("cannot write timestamp to file")
		return
	}
	return
}

func (c *Controller) updateLastTimestamp(timestamp string) (err error) {
	c.Logger.Printf("Updating last timestamp to: %s", timestamp)
	c.lastTimestamp = timestamp
	err = c.writeLastTimestamp()
	return
}
