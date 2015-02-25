package main

import (
	"golang.org/x/net/websocket"
)

type manager struct {
	c *websocket.Conn
	d chan bool
}

func socketHandler(conn *websocket.Conn) {
	for {
		// Send connection to matchChan.
		doneChan := make(chan bool)
		matchChan <- manager{c: conn, d: doneChan}

		ok := <-doneChan
		if !ok {
			return
		}
	}
}
