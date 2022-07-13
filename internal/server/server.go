package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Server struct {
	Rooms    map[string]Room
	Dispatch *chan Action
}

// Initialize a new server
func New() *Server {
	ch := make(chan Action)
	return &Server{
		Rooms:    make(map[string]Room),
		Dispatch: &ch,
	}
}

// Receive all incoming actions and map them to appropriate handlers
func (s *Server) Receive() {
	for {
		action := <-*s.Dispatch
		switch action.Type {
		case NAME_SET:
			s.SetName(action)
		case ROOM_JOIN:
			s.JoinRoom(action)
		case ROOM_LIST:
			s.ListRooms(action)
		case MESSAGE:
			s.Message(action)
		case ROOM_LEAVE:
			s.LeaveRoom(action)
		}
	}
}

// Creates an entry point for the server to work as a websocket handler
func (s *Server) Handler(upgrader websocket.Upgrader) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Connection error: %s", err.Error()))
		}
		fmt.Fprintf(os.Stdout, "New connection %s\n", conn.RemoteAddr())

		c := NewClient(conn, s)
		go c.Recieve()
	}

}

func (s *Server) SetName(a Action) {
	name, ok := a.Payload.(string)
	if !ok {
		a.Client.Error("Invalid name")
		return
	}
	a.Client.SetName(name)

}
func (s *Server) JoinRoom(a Action) {
	roomName, ok := a.Payload.(string)
	if !ok {
		a.Client.Error("Invalid room name")
		return
	}

	r, ok := s.Rooms[roomName]

	if !ok {
		r = *NewRoom(roomName)
		s.Rooms[roomName] = r
	}
	a.Client.JoinRoom(&r)

}

func (s *Server) ListRooms(a Action) {

	rooms := []string{}
	for _, r := range s.Rooms {
		rooms = append(rooms, r.Name)
	}

	b, err := json.Marshal(rooms)
	if err != nil {
		fmt.Printf("Unable to list rooms due to error: %s", err.Error())
		a.Client.Error("Unable to list rooms")
		return
	}

	ac := Action{
		Type:    ROOM_LIST,
		Payload: string(b),
	}
	a.Client.Conn.WriteJSON(ac)
}

func (s *Server) Message(a Action) {
	msg, ok := a.Payload.(string)
	if !ok {
		a.Client.Error("Invalid message")
		return
	}
	a.Client.Chat(msg)

}

func (s *Server) LeaveRoom(a Action) {
	a.Client.LeaveRoom()
}
