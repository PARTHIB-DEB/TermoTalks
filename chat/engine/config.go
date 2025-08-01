package Engine

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// HTTP to WS change
var connupgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for WebSocket connections
	},
}

// Variable for ClientPool and Broadcast Channel
var ClientPool = make(map[*websocket.Conn]bool) // Connected clients
var Broadcast = make(chan Message)              // Broadcast channel for messages

// Function to Upgrade HTTP connections to WS connections
func UpdateConnections(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	conn, err := connupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	ClientPool[conn] = true
	return conn
}
