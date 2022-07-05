package server

const (
	SET_NAME   = "SET_NAME"
	JOIN_ROOM  = "JOIN_ROOM"
	LIST_ROOMS = "LIST_ROOMS"
	ERROR      = "ERROR"
	MESSGE     = "MESSAGE"
	LEAVE_ROOM = "LEAVE_ROOM"
)

type Action struct {
	Client  *Client `json:"-"`
	Type    string  `json:"type"`
	Payload string  `json:"payload"`
}
