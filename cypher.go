package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
)

func encrypt(key string, b []byte) (encrypted []byte, err error) {

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

func decrypt(key string, b []byte) (decrypted []byte, err error) {

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

func CreateEncryptor(key string) func(b []byte) ([]byte, error) {
	return func(b []byte) ([]byte, error) {
		return encrypt(key, b)
	}
}

func CreateDecryptor(key string) func(b []byte) ([]byte, error) {
	return func(b []byte) ([]byte, error) {
		return decrypt(key, b)
	}
}

type EncryptedString string

func (s EncryptedString) MarshalJSON() ([]byte, error) {
	b, err := Encryptor([]byte(s))
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

	dec, err := Decryptor(b1)
	if err != nil {
		return err
	}
	*s = EncryptedString(dec)
	return nil
}
