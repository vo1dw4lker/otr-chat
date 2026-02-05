package core

type Room struct {
	id      string
	clients map[*Client]bool
}

func NewRoom(roomID string) *Room {
	return &Room{
		id:      roomID,
		clients: make(map[*Client]bool),
	}
}
