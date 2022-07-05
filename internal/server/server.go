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
			action.Client.SetName(action.Payload)
		case JOIN_ROOM:
			s.Join(action.Payload, action.Client)
		case LIST_ROOMS:
			s.ListRooms(action.Client)
		case MESSGE:
			action.Client.Broadcast(action.Payload)
		case LEAVE_ROOM:
			action.Client.LeaveRoom()
		default:
			action.Client.Error(fmt.Sprintf("Unrecognized action %s", action.Type))
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

func (s *Server) Join(roomName string, c *Client) {

	r, ok := s.Rooms[roomName]

	if !ok {
		r = *NewRoom(roomName)
		s.Rooms[roomName] = r
	}
	c.JoinRoom(&r)

}

func (s *Server) ListRooms(c *Client) {

	var rooms []string
	for _, r := range s.Rooms {
		rooms = append(rooms, r.Name)
	}

	b, err := json.Marshal(rooms)
	if err != nil {
		fmt.Printf("Unable to list rooms due to error: %s", err.Error())
		c.Error("Unable to list rooms")
		return
	}

	a := Action{
		Type:    LIST_ROOMS,
		Payload: string(b),
	}
	c.Conn.WriteJSON(a)
}
