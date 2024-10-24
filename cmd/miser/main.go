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
	// Initialization of cypher:
	miser.InitCypher(strings.Repeat("0123", 8))

	// Create repositories:
	ar := miser.CreateAccountRegistry()
	tr := miser.CreateTransactionRegistry()
	br := miser.CreateBalanceRegistry()
	cr := miser.CreateCurrencyRegistry()
	tg := miser.CreateTagRegistry()
	tm := miser.CreateTagsMapRegistry()

	// Create service:
	l := miser.CreateLedger(ar, br, tr, cr, tg, tm)

	n, err := ar.Load()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d accounts loaded, err: %v\n", n, err)
	fmt.Printf("Accounts: %#v\n", ar.List())

	n, err = tr.Load()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d transactions loaded, err: %v\n", n, err)
	fmt.Printf("Transactions: %#v\n", tr.List())

	n, err = br.Load()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d balances loaded, err: %v\n", n, err)
	fmt.Printf("Balances: %#v\n", br.List())
	//	fmt.Println("check balance:", miser.CheckBalance())

	// load tags:
	n, err = tg.Load()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d tags loaded, err: %v\n", n, err)
	fmt.Printf("Tags: %#v\n", tg.List())

	// load tags map:
	n, err = tm.Load()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d tags map loaded, err: %v\n", n, err)

	ac1, err := l.CreateAccount(
		"SMBC Trust Bank", miser.Asset, "Salary account", "JPY", time.Now(), 1555.13)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Aeon, err := l.CreateAccount(
		"AEON Supermarket", miser.Expense, "work bank account", "JPY", time.Now(), 0)
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

	ac1B := l.AccountAmount(ac1.ID)
	fmt.Printf("Balance of SMBC before transaction: %.2f\n", ac1B)

	t1, err := l.CreateTransaction(ac1.ID, Aeon.ID, time.Now(), 112.56, "私は店に行き、卵2kgと小麦粉を買いました。")
	if err != nil {
		fmt.Println("create transaction failure:", err)
		os.Exit(1)
	}
	fmt.Printf("Transaction: %#v\n", t1)

	ac1B = l.AccountAmount(ac1.ID)
	fmt.Printf("Balance of SMBC after transaction: %.2f\n", ac1B)

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
	fmt.Println("Amount:", l.AmountTransaction(t1))
	fmt.Println("Balances:", br)
	// fmt.Println("check balance:", miser.CheckBalance())

	n, err = ar.Save()
	fmt.Printf("%d new accounts saved, err: %v\n", n, err)
	n, err = tr.Save()
	fmt.Printf("%d new transactions saved, err: %v\n", n, err)
	n, err = br.Save()
	fmt.Printf("%d balances saved, err: %v\n", n, err)
	n, err = tg.Save()
	fmt.Printf("%d tags saved, err: %v\n", n, err)
	n, err = tm.Save()
	fmt.Printf("%d tags map saved, err: %v\n", n, err)
}
