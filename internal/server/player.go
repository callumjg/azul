package server

type Player struct {
	*Client
	Score uint8 `json:"score"`
	Board Board `json:"board"`
}
