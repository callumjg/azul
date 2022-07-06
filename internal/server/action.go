package server

type ActionType = uint8

const (
	UNDEFINED ActionType = iota
	SET_NAME
	JOIN_ROOM
	LIST_ROOMS
	ERROR
	MESSGE
	LEAVE_ROOM
)

func GetActionType(s string) ActionType {
	switch s {
	case "SET_NAME":
		return SET_NAME
	case "JOIN_ROOM":
		return JOIN_ROOM
	case "LIST_ROOMS":
		return LIST_ROOMS
	case "ERROR":
		return ERROR
	case "MESSGE":
		return MESSGE
	case "LEAVE_ROOM":
		return LEAVE_ROOM
	default:
		return UNDEFINED
	}
}

type Action struct {
	Client  *Client     `json:"-"`
	Type    ActionType  `json:"type"`
	Payload interface{} `json:"payload"`
}
