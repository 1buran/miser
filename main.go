package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var Encryptor, Decryptor func(b []byte) ([]byte, error)

func main() {
	Encryptor = CreateEncryptor(strings.Repeat("0123", 8))
	Decryptor = CreateDecryptor(strings.Repeat("0123", 8))

	fmt.Println(LoadAccounts(), "accounts loaded")
	fmt.Printf("Accounts: %#v\n", Accounts)

	fmt.Println(LoadTransactions(), "transactions loaded")
	fmt.Printf("Accounts: %#v\n", Transactions)

	fmt.Println(LoadBalances(), "balances loaded")
	fmt.Printf("Balances: %#v\n", Balances)
	fmt.Println("check balance:", CheckBalance())

	ac1, err := CreateAccount("SMBC Trust Bank", Asset, "Salary account", "JPY")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Aeon, err := CreateAccount("AEON Supermarket", Expense, "work bank account", "JPY")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b, err := json.Marshal(ac1)
	if err != nil {
		fmt.Println("marshal failure:", err)
		os.Exit(1)
	}

	fmt.Printf("\n%s\n", string(b))

	var ac2 Account
	if err := json.Unmarshal(b, &ac2); err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	fmt.Printf("\n%#v\n", ac2)

	t1, err := CreateTransation(ac1.ID, Aeon.ID, time.Now(), "112.56", "私は店に行き、卵2kgと小麦粉を買いました。")
	if err != nil {
		fmt.Println("create transaction failure:", err)
		os.Exit(1)
	}

	b, err = json.Marshal(t1)
	if err != nil {
		fmt.Println("marshal failure:", err)
		os.Exit(1)
	}
	fmt.Printf("\n%s\n", string(b))

	var t1e Transaction
	if err := json.Unmarshal(b, &t1e); err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	fmt.Printf("\n%#v\n", ac2)
	fmt.Printf("\n%#v\n", t1e)
	fmt.Println("Amount:", t1.Amount())
	fmt.Println("Balances:", Balances)
	fmt.Println("check balance:", CheckBalance())
	fmt.Println(SyncAccounts(), "new accounts saved")
	fmt.Println(SyncTransactions(), "new transactions saved")
	fmt.Println(SaveBalances(), "balances saved")
}
