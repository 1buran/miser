package miser

import (
	"time"
)

type NumericID int64

func CreateID() NumericID {
	return NumericID(time.Now().UnixNano())
}
