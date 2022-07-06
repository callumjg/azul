package server

type ActionType = string

const (
	ERROR       ActionType = "ERROR"
	MESSAGE     ActionType = "MESSAGE"
	ROOM_JOIN   ActionType = "ROOM_JOIN"
	ROOM_LEAVE  ActionType = "ROOM_LEAVE"
	ROOM_LIST   ActionType = "ROOM_LIST"
	ROOM_SET    ActionType = "ROOM_SET"
	NAME_SET    ActionType = "NAME_SET"
	USER_UPDATE ActionType = "USER_UPDATE"
)

type Action struct {
	Client  *Client     `json:"-"`
	Type    ActionType  `json:"type"`
	Payload interface{} `json:"payload"`
}

func SActionRoomJoin(c *Client) Action {
	return Action{
		Type:    ROOM_JOIN,
		Payload: c,
	}
}

func SActionRoomSet(r *Room) Action {
	return Action{
		Type:    ROOM_SET,
		Payload: r,
	}
}

func SActionRoomLeave(c *Client) Action {
	return Action{
		Type:    ROOM_LEAVE,
		Payload: c,
	}
}

func SActionUserUpdate(c *Client) Action {
	return Action{
		Type:    USER_UPDATE,
		Payload: c,
	}
}

func SActionMessage(msg ChatMessage) Action {
	return Action{
		Type:    MESSAGE,
		Payload: msg,
	}
}
