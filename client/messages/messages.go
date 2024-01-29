package messages

import (
	"clientenc/encryption"
	"clientenc/models"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
)

func ReadConn(conn net.Conn, key []byte, dst chan string) {
	buf := make([]byte, 1024)

	for {
		read, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Fatalln("Connection closed!")
			} else {
				log.Fatalln("Err reading:", err)
			}
		}
		packet := &models.Package{}

		err = proto.Unmarshal(buf[:read], packet)
		if err != nil {
			log.Println(err)
		}

		err = encryption.Decrypt(packet, key)
		if err != nil {
			log.Fatalln(err)
		}
		dst <- fmt.Sprintf("%s: %s", packet.Name, packet.Data)

	}
}

func ListenForMessages(conn net.Conn, src chan string, name string, key []byte) {
	for {
		buf := <-src
		sendMessage(name, buf, key, conn)
	}
}

func sendMessage(name, msg string, key []byte, dst io.Writer) {
	packet := &models.Package{
		Name: name,
		Data: []byte(msg),
	}

	err := encryption.Encrypt(packet, key)
	if err != nil {
		log.Fatalln(err)
	}

	bytes, err := proto.Marshal(packet)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = dst.Write(bytes)
	if err != nil {
		log.Fatalln(err)
	}
}
