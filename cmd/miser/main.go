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

	n, err := miser.LoadAccounts()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d accounts loaded, err: %v\n", n, err)
	fmt.Printf("Accounts: %#v\n", miser.Accounts.Items)

	n, err = miser.LoadTransactions()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d transactions loaded, err: %v\n", n, err)
	fmt.Printf("Accounts: %#v\n", miser.Transactions.Items)

	n, err = miser.LoadBalances()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d balances loaded, err: %v\n", n, err)
	fmt.Printf("Balances: %#v\n", miser.Balances.Items)
	//	fmt.Println("check balance:", miser.CheckBalance())

	n, err = miser.LoadTags()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d tags loaded, err: %v\n", n, err)
	fmt.Printf("Balances: %#v\n", miser.Tags.Items)

	n, err = miser.LoadTagsMap()
	fmt.Println(strings.Repeat("---", 40))
	fmt.Printf("%d tags map loaded, err: %v\n", n, err)
	fmt.Printf("Balances: %#v\n", miser.TagsMap.Items)

	ac1, err := miser.CreateAccount("SMBC Trust Bank", miser.Asset, "Salary account", "JPY", 1555.13)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Aeon, err := miser.CreateAccount("AEON Supermarket", miser.Expense, "work bank account", "JPY", 0)
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

	ac1B := miser.Balances.AccountValue(ac1.ID)
	fmt.Printf("Balance of SMBC before transaction: %.2f\n", ac1B)

	t1, err := miser.CreateTransation(ac1.ID, Aeon.ID, time.Now(), 112.56, "私は店に行き、卵2kgと小麦粉を買いました。")
	if err != nil {
		fmt.Println("create transaction failure:", err)
		os.Exit(1)
	}
	fmt.Printf("Transaction: %#v\n", t1)

	ac1B = miser.Balances.AccountValue(ac1.ID)
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
	fmt.Println("Amount:", t1.Amount())
	fmt.Println("Balances:", miser.Balances)
	// fmt.Println("check balance:", miser.CheckBalance())

	n, err = miser.SyncAccounts()
	fmt.Printf("%d new accounts saved, err: %v\n", n, err)
	n, err = miser.SyncTransactions()
	fmt.Printf("%d new transactions saved, err: %v\n", n, err)
	n, err = miser.SaveBalances()
	fmt.Printf("%d balances saved, err: %v\n", n, err)

	n, err = miser.SaveTags()
	fmt.Printf("%d tags saved, err: %v\n", n, err)
	n, err = miser.SaveTagsMap()
	fmt.Printf("%d tags map saved, err: %v\n", n, err)

}
