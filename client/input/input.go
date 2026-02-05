package input

import (
	"bufio"
	"clientenc/models"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
)

func EnterRoom(conn net.Conn) {
	buf := make([]byte, 24)
	for {
		room := Prompt("Enter room id: ")

		request := &models.RoomRequest{
			RoomId: room,
		}
		requestByte, err := proto.Marshal(request)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = conn.Write(requestByte)
		if err != nil {
			log.Fatalln(err)
		}

		read, err := conn.Read(buf)
		if err != nil {
			log.Fatalln(err)
		}

		responce := &models.RoomResponse{}
		err = proto.Unmarshal(buf[:read], responce)
		if err != nil {
			log.Fatalln(err)
		}

		switch responce.Status {
		case models.Status_RoomAssigned:
			fmt.Println("Room assigned")
			return

		case models.Status_RoomFull:
			fmt.Println("Room is Full")
		}
	}
}

func Prompt(text string) (userInput string) {
	fmt.Print(text)
	reader := bufio.NewReader(os.Stdin)
	userInput, _ = reader.ReadString('\n')

	userInput = strings.TrimSpace(userInput)
	return
}

// GetSocket prompts the user to enter a server IP address and port
// and returns the combined server address in the format "host:port".
func GetSocket() (sock string) {
	ip := Prompt("Enter server IP address: ")
	port := Prompt("Enter server port (default is 7575): ")

	if port == "" {
		port = "7575" // Default port
	}

	sock = fmt.Sprintf("%s:%s", ip, port)
	return
}
