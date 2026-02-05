package encryption

import (
	"clientenc/models"
	"crypto/mlkem"
	"crypto/rand"
	"errors"
	"net"

	"golang.org/x/crypto/chacha20poly1305"
)

const prefix = "$#"

func Encrypt(packet *models.Package, key []byte) error {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return err
	}

	nonce := make([]byte, aead.NonceSize())
	_, err = rand.Read(nonce)
	packet.Nonce = nonce
	if err != nil {
		return err
	}

	packet.Data = aead.Seal(nil, nonce, packet.Data, nil)

	return nil
}

func Decrypt(packet *models.Package, key []byte) error {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return err
	}

	packet.Data, err = aead.Open(nil, packet.Nonce, packet.Data, nil)
	if err != nil {
		return err
	}

	return nil
}

func MLKEMKeyExchange(conn net.Conn) ([]byte, error) {
	// when joined first: nothing
	// when joined second: ask first if rdy
	_, err := conn.Write([]byte(prefix + "rdy?"))
	if err != nil {
		return nil, err
	}

	// when joined first: wait for second
	// when joined second: read answer to rdy
	buf := make([]byte, 2048)
	read, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	var sharedKey []byte
	switch string(buf[:read]) {
	// joined second
	case prefix + "rdy?":
		// reply to rdy
		_, err = conn.Write([]byte(prefix + "yeah"))
		if err != nil {
			return nil, err
		}

		// read peer's encapsulation key
		read, err = conn.Read(buf)
		if err != nil {
			return nil, err
		}
		encapKey, err := mlkem.NewEncapsulationKey1024(buf[:read])
		if err != nil {
			return nil, err
		}

		// encapsulate and send ciphertext
		ciphertext, secret := encapKey.Encapsulate()
		_, err = conn.Write(ciphertext)
		if err != nil {
			return nil, err
		}

		// make and return shared key
		sharedKey = secret

	// joined first
	case prefix + "yeah":
		// generate decapsulation key
		decapKey, err := mlkem.GenerateKey1024()
		if err != nil {
			return nil, err
		}

		// send our encapsulation key
		_, err = conn.Write(decapKey.EncapsulationKey().Bytes())
		if err != nil {
			return nil, err
		}

		// read ciphertext from peer
		read, err = conn.Read(buf)
		if err != nil {
			return nil, err
		}

		// make and return shared key
		sharedKey, err = decapKey.Decapsulate(buf[:read])
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("unexpected handshake message")
	}
	return sharedKey, nil
}
