package miser

import (
	"testing"
)

func TestCreateAccount(t *testing.T) {

	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		acc, err := CreateAccount("Deposit", Asset, "deposit account", "USD", 0.00)
		if err != nil {
			t.Fatal(err)
		}

		if a := Accounts.Get(acc.ID); a == nil {
			t.Fatal("account not found in registry of accounts")
		}

	})

	// todo add cases with error
}
