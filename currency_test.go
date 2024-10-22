package miser

import (
	"testing"
)

func TestCurrency(t *testing.T) {
	t.Parallel()

	cr := CreateCurrencyRegistry()
	c := cr.Get("USD")
	if c == nil {
		t.Fatal("currecny not found, iso code: USD")
	}

	if c.Code != "USD" {
		t.Errorf("expected usd, got: %#v", c)
	}
}
