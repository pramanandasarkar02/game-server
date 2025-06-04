package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	// Send a message every second
	go func() {
		for {
			err := conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
			if err != nil {
				log.Println("Write error:", err)
				return
			}
			time.Sleep(time.Second)
		}
	}()

	// Read messages from the server
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}
		fmt.Printf("Received: %s\n", message)
	}
}