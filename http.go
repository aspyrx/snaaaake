package main

import (
	"golang.org/x/net/websocket"
	"net/http"
)

type handler struct{}

func staticHandler(writer http.ResponseWriter, request *http.Request) {}

func (r *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/socket/" {
		websocket.Handler(socketHandler).ServeHTTP(writer, request)
		return
	}

	// Everything else should be static.
	staticHandler(writer, request)
}
