package main

import (
	"net/http"
)

type game struct {
	players []*player
	board   [32][32]int
}

type player struct {
	direction  byte
	head, tail coordinate
	path       []coordinate
	m          manager
	dead       bool
}

type coordinate struct {
	x, y int
}

func initGame(players []*player) *game {
	g := &game{}
	g.players = players

	return g
}

func main() {
	http.ListenAndServe("0.0.0.0:8080", &handler{})
}
