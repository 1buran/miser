package miser

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
)

const (
	ACCOUNTS_FILE     = "miser.ac"
	TRANSACTIONS_FILE = "miser.tr"
	BALANCE_FILE      = "miser.bl"
)

type Entities interface {
	Account | Transaction | Balance
}

type Registry[E Entities] interface {
	AccountRegistry | TransactionRegistry | BalanceRegistry

	Add(e E) int
	SyncQueued() []E
}

func Save[E Entities, R Registry[E]](registry R, fpath string) (n int) {
	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for _, item := range registry.SyncQueued() {
		b, err := json.Marshal(item)
		if err != nil {
			log.Printf("marshaling failure: %s, data: %#v", err, item)
		}

		b = append(b, 10) // add new line at the end

		if _, err := f.Write(b); err != nil {
			log.Printf("write failure: %s", err)
		}
		n++
	}
	return
}

func Load[E Entities, R Registry[E]](registry R, fpath string) (n int) {
	f, err := os.Open(fpath)
	if err != nil {
		log.Println(err)
		return
	}

	defer f.Close()

	r := bufio.NewReader(f)

	for {
		b, err := r.ReadBytes(10)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		var e E
		json.Unmarshal(b, &e)
		n += registry.Add(e)
	}
	return
}

// Remove all entites from journal marked for deletion.
func CleanUp[E Entities](fpath string) int {
	// todo delete from journal lines with Deleted: true
	return 0
}
