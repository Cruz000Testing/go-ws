package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan int)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	go func() {
		for {
			cl := <-broadcast

			for client := range clients {
				client.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(cl)))
			}
		}
	}()

	http.ListenAndServe(":8080", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error updating connection", err)
		return
	}

	defer ws.Close()

	clients[ws] = true

	broadcast <- len(clients)

	for {
		_, _, err := ws.ReadMessage()

		if err != nil {
			ws.Close()
			delete(clients, ws)
			break
		}
	}
}
