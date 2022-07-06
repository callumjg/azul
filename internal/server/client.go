package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

type Client struct {
	ID     uuid.UUID
	Name   string
	Conn   *websocket.Conn
	Room   *Room
	Server *Server
}

func NewClient(conn *websocket.Conn, s *Server) *Client {

	return &Client{
		ID:     uuid.New(),
		Name:   "Anon",
		Conn:   conn,
		Server: s,
	}
}

func (c *Client) Recieve() {
	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Printf("Closing connection %s\n", c.Conn.RemoteAddr().String())
			c.LeaveRoom()
			c.Conn.Close()
			break
		}
		var raw struct {
			Type    string      `json:"type"`
			Payload interface{} `json:"payload"`
		}

		err = json.Unmarshal(m, &raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read message from %s: %s\n", c.Conn.RemoteAddr().String(), err.Error())
			c.Error("Invalid message received")
			continue
		}
		a := Action{
			Client:  c,
			Type:    GetActionType(raw.Type),
			Payload: raw.Payload,
		}

		*c.Server.Dispatch <- a
	}
}

func (c *Client) SetName(name string) {
	oldName := c.Name
	c.Name = name
	c.Broadcast(fmt.Sprintf("%s changed their name to %s", oldName, name))

}

func (c *Client) Error(msg string) {
	a := Action{
		Type:    ERROR,
		Payload: msg,
	}
	c.Conn.WriteJSON(a)
}

func (c *Client) JoinRoom(r *Room) {
	c.LeaveRoom()
	r.Clients[c.ID] = c
	c.Room = r

	c.Broadcast(fmt.Sprintf("%s has joined the room\n", c.Name))
	c.Message(fmt.Sprintf("Welcome to %s\n", r.Name))
}

func (c *Client) LeaveRoom() {
	if c.Room != nil {
		delete(c.Room.Clients, c.ID)
		c.Broadcast(fmt.Sprintf("%s has let the room\n", c.Name))

		if len(c.Room.Clients) == 0 {
			delete(c.Server.Rooms, c.Room.Name)
		}
		c.Room = nil
	}
}

// Send message to all other members of a room
func (c *Client) Broadcast(msg string) {
	if c.Room == nil {
		return
	}
	for id, m := range c.Room.Clients {
		if id != c.ID {
			m.Message(msg)
		}
	}

}

// Send message to client
func (c *Client) Message(msg string) {
	a := Action{
		Type:    MESSGE,
		Payload: msg,
	}
	c.Conn.WriteJSON(a)
}
