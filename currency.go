package miser

import (
	_ "embed"
	"encoding/json"
)

//go:embed currency.json
var currencyJsonContent []byte

type Currency struct{ Code, Name, Sign string }

type CurrencyRegistry map[string]Currency

func (cr CurrencyRegistry) Get(code string) *Currency {
	c, ok := cr[code]
	if !ok {
		return nil
	}
	return &c
}

func CreateCurrencyRegistry() *CurrencyRegistry {
	cr := make(CurrencyRegistry)
	err := json.Unmarshal(currencyJsonContent, &cr)
	if err != nil {
		panic(err)
	}
	return &cr
}
