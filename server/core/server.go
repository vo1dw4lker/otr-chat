package core

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	clients       map[*Client]bool
	rooms         map[string]*Room
	addClient     chan *Client
	removeClient  chan *Client
	clientMessage chan *ClientMessage
}

func NewServer() *Server {
	return &Server{
		clients:       make(map[*Client]bool),
		rooms:         make(map[string]*Room),
		addClient:     make(chan *Client),
		removeClient:  make(chan *Client),
		clientMessage: make(chan *ClientMessage),
	}
}

func (server *Server) Start(port int) {
	go server.handleChannels()

	listener, err := net.Listen("tcp",
		fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln(err)
	}

	defer listener.Close()
	log.Printf("Server started on port %d", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			break
		}

		client := NewClient(server, conn)
		go server.handleClient(client)
	}
}
func (server *Server) handleChannels() {
	for {
		select {
		case client := <-server.addClient:
			server.clients[client] = true
			log.Printf("New client connected: %v", client.conn.RemoteAddr().String())

		case client := <-server.removeClient:
			delete(server.clients, client)

			if client.room != nil {
				delete(client.room.clients, client)

				if len(client.room.clients) != 0 {
					for v := range client.room.clients {
						v.conn.Close()
					}
				}

				delete(server.rooms, client.room.id)
			}

			log.Printf("Client disconnected: %v", client.conn.RemoteAddr().String())

		case msg := <-server.clientMessage:
			for client := range msg.client.room.clients {
				if client != msg.client {
					_, err := client.conn.Write(msg.message)
					if err != nil {
						handleError(server, client, WriteError, err)
					}
				}
			}
		}
	}
}

func (server *Server) handleClient(client *Client) {
	buffer := make([]byte, 1024)

	client.assignToRoom(server)
	for {
		n, err := client.conn.Read(buffer)
		if err != nil {
			handleError(server, client, ReadError, err)
			break
		}

		message := buffer[:n]

		server.clientMessage <- NewClientMessage(client, message)
	}
}
