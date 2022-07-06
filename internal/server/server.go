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

func New() *Server {
	ch := make(chan Action)
	return &Server{
		Rooms:    make(map[string]Room),
		Dispatch: &ch,
	}
}

func (s *Server) Receive() {
	for {
		action := <-*s.Dispatch
		switch action.Type {
		case SET_NAME:
			s.SetName(action)
		case JOIN_ROOM:
			s.JoinRoom(action)
		case LIST_ROOMS:
			s.ListRooms(action)
		case MESSGE:
			s.Message(action)
		case LEAVE_ROOM:
			s.LeaveRoom(action)
		default:
			action.Client.Error("Unrecognized action")
		}
	}
}

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

	var rooms []string
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
		Type:    LIST_ROOMS,
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
	a.Client.Broadcast(msg)

}

func (s *Server) LeaveRoom(a Action) {
	a.Client.LeaveRoom()
}
