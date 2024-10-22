package miser

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
)

// Cypher provides ecrypt and decrypt methods.
type Cypher struct {
	encrypt func(b []byte) ([]byte, error)
	decrypt func(b []byte) ([]byte, error)
}

var cypher Cypher

// Init cypher service for a given key.
func InitCypher(key string) {
	cypher = Cypher{
		decrypt: func(b []byte) ([]byte, error) {
			return decryptor(key, b)
		},
		encrypt: func(b []byte) ([]byte, error) {
			return encryptor(key, b)
		},
	}
}

func encryptor(key string, b []byte) (encrypted []byte, err error) {

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted = aesgcm.Seal(nonce, nonce, b, nil)
	return
}

func decryptor(key string, b []byte) (decrypted []byte, err error) {

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()

	nonce, encrypted := b[:nonceSize], b[nonceSize:]
	if decrypted, err = aesgcm.Open(nil, []byte(nonce), []byte(encrypted), nil); err != nil {
		return nil, err
	}

	return
}

type EncryptedString string

func (s EncryptedString) MarshalJSON() ([]byte, error) {
	b, err := cypher.encrypt([]byte(s))
	if err != nil {
		return nil, err
	}
	return json.Marshal(b)
}

func (s *EncryptedString) UnmarshalJSON(b []byte) error {
	var b1 []byte
	if err := json.Unmarshal(b, &b1); err != nil {
		return err
	}

	dec, err := cypher.decrypt(b1)
	if err != nil {
		return err
	}
	*s = EncryptedString(dec)
	return nil
}
