package main

import (
	"net/http"

	Engine "github.com/PARTHIB-DEB/TermoTalks/chat/engine"
)

func main() {
	// Initialize the chat engine
	Engine.Apiroutes()

	// Start the HTTP server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
