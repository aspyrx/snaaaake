package main

import (
	"fmt"
	"math/rand"
	"time"
)

type delta struct {
	x, y, value int
}

const MAP_WIDTH int = 32
const MAP_HEIGHT int = 32

func (g *game) randFood() coordinate {
	for {
		x := rand.Intn(MAP_WIDTH)
		y := rand.Intn(MAP_HEIGHT)
		if g.board[x][y] == '-' {
			return coordinate{x: x, y: y}
		}
	}
}

func gameHandle(g *game) {
	remaining := len(g.players)

	for x := 0; x < MAP_WIDTH; x++ {
		for y := 0; y < MAP_HEIGHT; y++ {
			g.board[x][y] = '-'
		}
	}

	for key, value := range g.players {
		// Set the default direction to the right.
		value.direction = 'l'

		// Let the users know we are ready.
		_, err := value.m.c.Write([]byte("start " + string('0'+key)))
		if err != nil {
			// Send false to doneChan to signal connection failure.
			value.m.d <- false
			// Remove user and snake from scene.
		}
	}

	// Read from user connections.
	for key, value := range g.players {
		go func(key int, value *player) {
			for {
				msg := make([]byte, 64)
				_, err := value.m.c.Read(msg)
				if err != nil {
					return
				}
				if len(msg) < 5 {
					return
				}
				if msg[4] == 'h' && value.direction != 'l' {
					value.direction = 'h'
				}
				if msg[4] == 'j' && value.direction != 'k' {
					value.direction = 'j'
				}
				if msg[4] == 'k' && value.direction != 'j' {
					value.direction = 'k'
				}
				if msg[4] == 'l' && value.direction != 'h' {
					value.direction = 'l'
				}
				switch msg[4] {
				case 'h', 'j', 'k', 'l':
					continue
				default:
					return
				}
			}
		}(key, value)
	}

	dlist := make([]delta, 0)

	rF := g.randFood()
	dlist = append(dlist, delta{x: rF.x, y: rF.y, value: 'x'})
	g.board[rF.x][rF.y] = 'x'

	for key, value := range g.players {
		value.tail = coordinate{0, key * 5}
		value.head = coordinate{8, key * 5}
		for i := 0; i <= 8; i++ {
			value.path = append(value.path, coordinate{i, key * 5})
			g.board[i][key*5] = '0' + key
			dlist = append(dlist, delta{x: i, y: key * 5, value: '0' + key})
		}
	}

	for {
		if remaining == 1 {
			winner := 4
			for key, value := range g.players {
				if !value.dead {
					winner = key
					break
				}
			}

			for _, value := range g.players {
				value.m.c.Write([]byte("end " + string('0'+winner)))
				value.m.c.Close()
			}
			return
		}

		// Construct the redraw string
		redraw := ""
		for _, value := range dlist {
			redraw += fmt.Sprintf("%d,%d,%s|", value.x, value.y, string(value.value))
		}

		if len(redraw) < 1 {
			// Lost connection at some point, end the game
			for _, value := range g.players {
				value.m.c.Write([]byte("end 4"))
				value.m.c.Close()
			}
			return
		}

		redraw = redraw[:len(redraw)-1]

		for _, value := range g.players {
			// Let the users know we are ready.
			_, err := value.m.c.Write([]byte("redraw " + redraw))
			if err != nil {
				// Send false to doneChan to signal connection failure.
				value.m.d <- false
				// Remove user and snake from scene.
				for _, coor := range value.path {
					g.board[coor.x][coor.y] = '-'
					//					fmt.Println(coor)
					dlist = append(dlist, delta{x: coor.x, y: coor.y, value: '-'})
					value.dead = true
					remaining--
				}
			}
		}

		dlist = []delta{}

		time.Sleep(100 * time.Millisecond)

		for key, value := range g.players {
			if value.dead {
				continue
			}

			oldTail := value.path[0]

			// Remove the tail
			g.board[value.tail.x][value.tail.y] = '-'
			dlist = append(dlist, delta{x: value.tail.x, y: value.tail.y, value: '-'})

			// Figure out the new tail
			value.tail = value.path[1]

			// Remove the old tail
			value.path = value.path[1:]

			// Determine the old head
			oldC := value.path[len(value.path)-1]

			// Calculate the new head
			newC := coordinate{}
			switch value.direction {
			case 'h':
				newC = coordinate{(oldC.x + MAP_WIDTH - 1) % MAP_WIDTH, oldC.y}
			case 'j':
				newC = coordinate{oldC.x, (oldC.y + 1) % MAP_HEIGHT}
			case 'k':
				newC = coordinate{oldC.x, (oldC.y + MAP_HEIGHT - 1) % MAP_HEIGHT}
			case 'l':
				newC = coordinate{(oldC.x + 1) % MAP_WIDTH, oldC.y}
			}

			// COLLISION! EXIT.
			switch g.board[newC.x][newC.y] {
			case '-':
				// Add the new head
				value.path = append(value.path, newC)
				g.board[newC.x][newC.y] = '0' + key
				dlist = append(dlist, delta{x: newC.x, y: newC.y, value: '0' + key})
			case 'x':
				value.tail = oldTail
				value.path = append([]coordinate{oldTail}, value.path...)
				dlist = append(dlist, delta{x: oldTail.x, y: oldTail.y, value: '0' + key})

				rF := g.randFood()
				dlist = append(dlist, delta{x: newC.x, y: newC.y, value: '0' + key})
				g.board[newC.x][newC.y] = '0' + key
				dlist = append(dlist, delta{x: rF.x, y: rF.y, value: 'x'})
				value.path = append(value.path, newC)
				g.board[rF.x][rF.y] = 'x'
			default:
				value.m.c.Write([]byte("death " + string('0'+key)))

				// Remove the user
				for _, coor := range value.path {
					g.board[coor.x][coor.y] = '-'
					//					fmt.Println(coor)
					dlist = append(dlist, delta{x: coor.x, y: coor.y, value: '-'})
					value.dead = true
				}
				remaining--
			}
		}
	}
}
