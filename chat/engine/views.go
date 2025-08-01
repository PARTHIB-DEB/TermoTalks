package Engine

import (
	"fmt"
	"log"
	"net/http"
)

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for sending messages via WebSocket
	// This function will handle the WebSocket connection and message
	wsconn := UpdateConnections(w, r)
	for {
		var msg Message
		err := wsconn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			delete(ClientPool, wsconn)
			return
		}
		Broadcast <- msg
		fmt.Println("Message received:", msg.content)
	}
}

func ReceiveMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for receiving messages via WebSocket
	// This function will handle the WebSocket connection and message receiving
	for {
		msg := <-Broadcast

		for client := range ClientPool {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(ClientPool, client)
			}
		}
	}
}
