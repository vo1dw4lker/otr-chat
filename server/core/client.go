package core

import (
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"serverenc/models"
)

type Client struct {
	conn net.Conn
	room *Room
}

type ClientMessage struct {
	client  *Client
	message []byte
}

func NewClient(server *Server, conn net.Conn) *Client {
	client := &Client{conn: conn}
	server.addClient <- client

	log.Printf("New client connected: %s", conn.RemoteAddr().String())
	return client
}

func NewClientMessage(client *Client, message []byte) *ClientMessage {
	return &ClientMessage{
		client:  client,
		message: message,
	}
}

// Receive room ID from a client and assign it to the client
// Reject if number of clients in the room is greater than 2
func (client *Client) assignToRoom(server *Server) {
	buf := make([]byte, 256)

	// Loop until correct (not full) room entered by user
	for {
		n, err := client.conn.Read(buf)
		if err != nil {
			handleError(server, client, ReadError, err)
			break
		}

		request := &models.RoomRequest{}
		err = proto.Unmarshal(buf[:n], request)
		if err != nil {
			handleError(server, client, InternalError, err)
		}

		if server.rooms[request.RoomId] == nil {
			server.rooms[request.RoomId] = NewRoom(request.RoomId)
		}

		if len(server.rooms[request.RoomId].clients) < 2 {
			client.room = server.rooms[request.RoomId]
			client.room.clients[client] = true

			response := &models.RoomResponse{
				Status: models.Status_RoomAssigned,
			}

			responseByte, err := proto.Marshal(response)
			if err != nil {
				handleError(server, client, InternalError, err)
			}

			// Respond with RoomAssigned status
			_, err = client.conn.Write(responseByte)
			if err != nil {
				handleError(server, client, WriteError, err)
			}
			break
		}

		response := &models.RoomResponse{
			Status: models.Status_RoomFull,
		}

		responseByte, err := proto.Marshal(response)
		if err != nil {
			handleError(server, client, InternalError, err)
		}

		// Respond with RoomFull status
		_, err = client.conn.Write(responseByte)
		if err != nil {
			handleError(server, client, WriteError, err)
			break
		}
	}
}
