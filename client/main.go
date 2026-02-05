package main

import (
	"clientenc/encryption"
	"clientenc/input"
	"clientenc/messages"
	"clientenc/ui"
	"log"
	"net"
)

func main() {
	sock := input.GetSocket()

	conn, err := net.Dial("tcp", sock)
	if err != nil {
		log.Fatalln("Unable to connect:", err)
	}
	defer conn.Close()

	name := input.Prompt("Enter your name: ")
	input.EnterRoom(conn)

	// key exchange
	key, err := encryption.MLKEMKeyExchange(conn)
	if err != nil {
		log.Fatalln(err)
	}

	rcvdMsg := make(chan string)
	sendMsg := make(chan string)

	go messages.ReadConn(conn, key, rcvdMsg)

	go messages.ListenForMessages(conn, sendMsg, name, key)

	if err := ui.RunChatUI(rcvdMsg, sendMsg); err != nil {
		log.Fatalln(err)
	}
}
