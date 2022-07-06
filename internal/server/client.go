package server

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"github.com/google/uuid"
)

type Client struct {
	ID     uuid.UUID       `json:"id"`
	Name   string          `json:"name"`
	Conn   *websocket.Conn `json:"-"`
	Room   *Room           `json:"-"`
	Server *Server         `json:"-"`
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
		var a Action

		err = json.Unmarshal(m, &a)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read message from %s: %s\n", c.Conn.RemoteAddr().String(), err.Error())
			c.Error("Invalid message received")
			continue
		}
		a.Client = c
		*c.Server.Dispatch <- a
	}
}

func (c *Client) SetName(name string) {
	c.Name = name
	c.BroadcastAction(SActionUserUpdate(c))
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
	c.BroadcastAction(SActionRoomJoin(c))
	c.Action(SActionRoomSet(r))
}

func (c *Client) LeaveRoom() {
	if c.Room != nil {
		delete(c.Room.Clients, c.ID)
		c.BroadcastAction(SActionRoomLeave(c))
		if len(c.Room.Clients) == 0 {
			delete(c.Server.Rooms, c.Room.Name)
		}
		c.Room = nil
	}
}

// Send message to all other members of a room
func (c *Client) Chat(msg string) {
	if c.Room == nil {
		return
	}

	chatMsg := ChatMessage{
		UID:     c.ID,
		Time:    time.Now().UnixMilli(),
		Message: msg,
	}

	c.Room.BroadcastChat(chatMsg)

}

func (c *Client) Action(a Action) {
	c.Conn.WriteJSON(a)
}

func (c *Client) BroadcastAction(a Action) {
	if c.Room == nil {
		return
	}
	for id, cl := range c.Room.Clients {
		if id == c.ID {
			continue
		}
		cl.Conn.WriteJSON(a)
	}
}
