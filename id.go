package miser

import (
	"crypto/rand"
	"fmt"
)

const RANDOM_BYTES_LENGTH = 10

type ID string

func CreateID() ID {
	b := make([]byte, RANDOM_BYTES_LENGTH)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return ID(fmt.Sprintf("%x", b))
}
