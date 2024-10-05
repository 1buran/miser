package main

import (
	"encoding/json"
	"fmt"
	"github.com/1buran/miser"
	"os"
	"strings"
	"time"
)

func main() {
	miser.Encryptor = miser.CreateEncryptor(strings.Repeat("0123", 8))
	miser.Decryptor = miser.CreateDecryptor(strings.Repeat("0123", 8))

	fmt.Println(miser.LoadAccounts(), "accounts loaded")
	fmt.Printf("Accounts: %#v\n", miser.Accounts)

	fmt.Println(miser.LoadTransactions(), "transactions loaded")
	fmt.Printf("Accounts: %#v\n", miser.Transactions)

	fmt.Println(miser.LoadBalances(), "balances loaded")
	fmt.Printf("Balances: %#v\n", miser.Balances)
	fmt.Println("check balance:", miser.CheckBalance())

	ac1, err := miser.CreateAccount("SMBC Trust Bank", miser.Asset, "Salary account", "JPY")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Aeon, err := miser.CreateAccount("AEON Supermarket", miser.Expense, "work bank account", "JPY")
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

	var ac2 miser.Account
	if err := json.Unmarshal(b, &ac2); err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	fmt.Printf("\n%#v\n", ac2)

	t1, err := miser.CreateTransation(ac1.ID, Aeon.ID, time.Now(), "112.56", "私は店に行き、卵2kgと小麦粉を買いました。")
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

	var t1e miser.Transaction
	if err := json.Unmarshal(b, &t1e); err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	fmt.Printf("\n%#v\n", ac2)
	fmt.Printf("\n%#v\n", t1e)
	fmt.Println("Amount:", t1.Amount())
	fmt.Println("Balances:", miser.Balances)
	fmt.Println("check balance:", miser.CheckBalance())
	fmt.Println(miser.SyncAccounts(), "new accounts saved")
	fmt.Println(miser.SyncTransactions(), "new transactions saved")
	fmt.Println(miser.SaveBalances(), "balances saved")
}
