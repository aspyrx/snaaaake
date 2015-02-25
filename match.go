package main

import (
	"sync"
	"time"
)

var matchChan chan manager
var matchLock *sync.Mutex
var matchKey = 0

func matchHandler() {
	g := &game{}
	i := 0
	for {
		// Listen on channel for connection.
		m := <-matchChan

		// Obtain lock to avoid goroutine screwing with matches.
		matchLock.Lock()

		// Add player to the game
		g.players = append(g.players, &player{m: m})

		// Increment obviously
		i++

		// If i == 1, this is the first player. We must wait for additional
		// players.
		if i == 1 {
			// Unlock and proceed to next connection.
			matchLock.Unlock()
			continue
		}

		// If i == 2, there are two existing people. Give additional people 10
		// seconds to join current game.
		if i == 2 {
			go func(localKey int) {
				// Wait 10 seconds for additional players.
				time.Sleep(10*time.Second)

				// We need a mutex on reading and writing to matchKey.
				matchLock.Lock()
				defer matchLock.Unlock()

				// If key is different from previous key, the previous game had
				// reached 4 people has started. Exit.
				if localKey != matchKey {
					return
				}

				// We have either 2 or 3 players. Start the game.
				go gameHandle(g)

				// Reset variables. Increase the matchKey.
				g = &game{}
				i = 0
				matchKey++
			}(matchKey)

			// Unlock and proceed to next connection.
			matchLock.Unlock()
			continue
		}

		// If i == 3, nothing needs to be done. Either it will be caught in the
		// goroutine above, or will hit 4 players below.
		if i == 3 {
			// Unlock and proceed to next connection.
			matchLock.Unlock()
			continue
		}

		// Full game with 4 players. Start the game.
		go gameHandle(g)

		// Reset variables. Increase the matchKey.
		g = &game{}
		i = 0
		matchKey++

		// Unlock and proceed to next connection.
		matchLock.Unlock()
	}
}

func init() {
	matchChan = make(chan manager)
	matchLock = &sync.Mutex{}
	go matchHandler()
}
