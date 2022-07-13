package server

import (
	"github.com/google/uuid"
)

type ChatMessage struct {
	UID     uuid.UUID `json:"uid"`
	Time    int64     `json:"time"`
	Message string    `json:"message"`
}

type Room struct {
	Name    string                `json:"name"`
	Clients map[uuid.UUID]*Client `json:"clients"`
	Chat    *[]ChatMessage        `json:"chat"`
	Game    *Game                 `json:"game"`
}

// Initialize a new room
func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		Clients: make(map[uuid.UUID]*Client),
		Chat:    &[]ChatMessage{},
	}
}

// Broadcast message to all members of a room
func (r *Room) BroadcastChat(msg ChatMessage) {
	*r.Chat = append(*r.Chat, msg)
	for _, c := range r.Clients {
		c.Conn.WriteJSON(SActionMessage(msg))
	}
}

// Broadcast action to all members of a room
func (r *Room) BroadcastAction(a Action) {
	for _, c := range r.Clients {
		c.Conn.WriteJSON(a)
	}
}
