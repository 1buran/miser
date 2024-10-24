package miser

import (
	"testing"
	"time"
)

func TestCreateAccount(t *testing.T) {

	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		// Create repositories:

		ar := CreateAccountRegistry()
		tr := CreateTransactionRegistry()
		br := CreateBalanceRegistry()
		cr := CreateCurrencyRegistry()
		tg := CreateTagRegistry()
		tm := CreateTagsMapRegistry()

		// Create service:
		l := CreateLedger(ar, br, tr, cr, tg, tm)

		acc, err := l.CreateAccount("Deposit", Asset, "deposit account", "USD", time.Now(), 0.00)
		if err != nil {
			t.Fatal(err)
		}

		if a := ar.Get(acc.ID); a == nil {
			t.Fatal("account not found in registry of accounts")
		}

	})

	// todo add cases with error
}
