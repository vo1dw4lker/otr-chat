package core

import "log"

type ErrType string

const (
	ReadError     ErrType = "Error reading from client %s: %s"
	WriteError            = "Error writing to client %s: %s"
	InternalError         = "Internal error, serving %s: %s"
)

func handleError(server *Server, client *Client, errType ErrType, err error) {
	log.Printf(string(errType), client.conn.RemoteAddr(), err)
	server.removeClient <- client
}
