package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/callumjg/azul/internal/server"
	"github.com/gorilla/websocket"
)

var connections map[string]*websocket.Conn

func receive(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for addr, c := range connections {
			if addr == conn.RemoteAddr().String() {
				continue
			}
			c.WriteMessage(messageType, p)
		}
	}
}

func main() {

	s := server.New()
	go s.Receive()

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", s.Handler(upgrader))
	port := ":8888"

	fmt.Printf("Listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Panic(err)
	}

}
