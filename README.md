# OTR Chat

OTR Chat is a secure communication project that uses OTR-like protocol for establishing encrypted conversations between two clients.

## Table of Contents

- [About](#about)
- [Usage](#usage)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## About

OTR Chat is a secure communication project designed to establish private conversations between two clients. The project consists of two main components: the server and the client.

### Server Functionality

The server manages rooms, assings clients to them and relays their messages. Even if server's communications are monitored, no one should be able to read the messages, as only the clients have the decryption keys.

### Client Security Measures

Client uses the X25509 algorithm to securely establish a shared key for communication. Then, every message is encrypted with chacha20-poly1305. This algorithm not only encrypts a string, but also verifies its integrity. This guarantees end-to-end encryption for every message exchanged between clients.

## Usage

1. Install [Go](https://go.dev) on your machine.
2. Clone or download the repository.
3. Within the *client* and *server* directories, run `go build .` to generate binaries. 

After building, you can execute the generated binaries to start the client and server applications.

## License

Licensed under the MIT License. See [LICENSE.txt](LICENSE.txt) for details.

## Acknowledgments

A big thanks to contributors of the [tview](https://github.com/rivo/tview) library.