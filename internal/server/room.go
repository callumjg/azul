package server

import (
	"github.com/google/uuid"
)

type Room struct {
	Name    string                `json:"name"`
	Clients map[uuid.UUID]*Client `json:"-"`
}

func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		Clients: make(map[uuid.UUID]*Client),
	}
}

// Send message to all members of a room
func (r *Room) Broadcast(msg string) {
	for _, c := range r.Clients {
		c.Message(msg)
	}
}
