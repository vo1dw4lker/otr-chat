package encryption

import (
	"clientenc/models"
	"crypto/ecdh"
	"crypto/rand"
	"golang.org/x/crypto/chacha20poly1305"
	"net"
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

func ECDHKeyExchange(conn net.Conn) ([]byte, error) {
	privKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	pubKey := privKey.PublicKey()

	// when joined first: nothing
	// when joined second: ask first if rdy
	_, err = conn.Write([]byte(prefix + "rdy?"))
	if err != nil {
		return nil, err
	}

	// when joined first: wait for second
	// when joined second: read answer to rdy
	buf := make([]byte, 256)
	read, err := conn.Read(buf)

	var sharedKey []byte
	switch string(buf[:read]) {
	// joined second
	case prefix + "rdy?":
		// reply to rdy
		_, err = conn.Write([]byte(prefix + "yeah"))
		if err != nil {
			return nil, err
		}

		// read peer's public key
		read, err = conn.Read(buf)
		peerPubKey, err := ecdh.X25519().NewPublicKey(buf[:read])
		if err != nil {
			return nil, err
		}

		// send our public key
		_, err = conn.Write(pubKey.Bytes())
		if err != nil {
			return nil, err
		}

		// make and return shared key
		sharedKey, err = privKey.ECDH(peerPubKey)
		if err != nil {
			return nil, err
		}

	// joined first
	case prefix + "yeah":
		// send our public key
		_, err = conn.Write(pubKey.Bytes())
		if err != nil {
			return nil, err
		}

		read, err = conn.Read(buf)
		peerPubKey, err := ecdh.X25519().NewPublicKey(buf[:read])
		if err != nil {
			return nil, err
		}

		// make and return shared key
		sharedKey, err = privKey.ECDH(peerPubKey)
		if err != nil {
			return nil, err
		}

	}
	return sharedKey, nil
}
