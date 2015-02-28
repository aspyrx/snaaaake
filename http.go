package main

import (
	"golang.org/x/net/websocket"
	"net/http"
)

type handler struct{}

func (r *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/socket/":
		websocket.Handler(socketHandler).ServeHTTP(writer, request)
		return
	case "/game.js":
		http.ServeFile(writer, request, "src/github.com/aspyrx/snaaaake/game.js")
	case "/style.css":
		http.ServeFile(writer, request, "src/github.com/aspyrx/snaaaake/style.css")
	default:
		http.ServeFile(writer, request, "src/github.com/aspyrx/snaaaake/index.html")
	}
}
