package main

import (
	"encoding/json"
	"strings"
	"testing"
)

type TestData struct {
	Regular string
	Secret  EncryptedString
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	t.Run("encrypt", func(t *testing.T) {
		b, err := encrypt(strings.Repeat("secret k", 4), []byte("some text"))
		if err != nil {
			t.Fatal(err)
		}

		s := (string(b))
		if strings.Contains(s, "some text") {
			t.Fatal("encryption failed")
		}
	})

	t.Run("decrypt", func(t *testing.T) {
		b, _ := encrypt(strings.Repeat("secret k", 4), []byte("some text"))
		dec, err := decrypt(strings.Repeat("secret k", 4), b)
		if err != nil {
			t.Fatal(err)
		}

		if string(dec) != "some text" {
			t.Fatal("decryption failed")
		}
	})

	t.Run("wrong key", func(t *testing.T) {
		b, _ := encrypt(strings.Repeat("secret k", 4), []byte("some text"))
		_, err := decrypt(strings.Repeat("sAcrAt A", 4), b)
		if err == nil {
			t.Error("error expected, nil found")
		} else {
			t.Log(err)
		}
	})
}

func TestEncryptedString(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		Encryptor = CreateEncryptor(strings.Repeat("0123", 8))
		Decryptor = CreateDecryptor(strings.Repeat("0123", 8))

		data := TestData{Regular: "hello", Secret: "something hidden"}

		b, err := json.Marshal(data)
		s := string(b)

		t.Logf("%#v", s)

		if err != nil {
			t.Fatal(err)
		}

		if strings.Contains(s, "something hidden") {
			t.Fatal("encryption failed")
		}
	})

	t.Run("JSON unmarshaling", func(t *testing.T) {
		Encryptor = CreateEncryptor(strings.Repeat("0123", 8))
		Decryptor = CreateDecryptor(strings.Repeat("0123", 8))

		data := TestData{Regular: "hello", Secret: "something hidden"}

		b, _ := json.Marshal(data)

		data1 := TestData{}

		if err := json.Unmarshal(b, &data1); err != nil {
			t.Fatal(err)
		}

		t.Logf("%#v", data1)

		if data1.Secret != "something hidden" {
			t.Errorf("Got unexpected string: %s", data1.Secret)
		}

	})

	t.Run("JSON unmarshaling with wrong key", func(t *testing.T) {
		Encryptor = CreateEncryptor(strings.Repeat("0123", 8))
		Decryptor = CreateDecryptor(strings.Repeat("abcd", 8))

		data := TestData{Regular: "hello", Secret: "something hidden"}

		b, _ := json.Marshal(data)

		data2 := TestData{}

		if err := json.Unmarshal(b, &data2); err == nil {
			t.Fatal("error auth failed expected")
		}

		t.Logf("encypted string shoud be empty: %#v", data2)
	})
}
