package Engine

import "net/http"

func Apiroutes() {
	http.HandleFunc("/send/ws", SendMessageHandler)
	http.HandleFunc("/receive/ws", ReceiveMessageHandler)
}
