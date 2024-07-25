package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"reflect"
)

const (
	ACCOUNTS_FILE     = "miser.ac"
	TRANSACTIONS_FILE = "miser.tr"
	BALANCE_FILE      = "miser.bl"
)

type Entities interface {
	Account | Transaction | Balance
}

func Save[E Entities](items map[NumericID]struct{}, registry map[NumericID]E, fpath string) (n int) {
	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for id := range items {
		b, err := json.Marshal(registry[id])
		if err != nil {
			log.Printf("marshaling failure: %s, data: %#v", err, registry[id])
		}

		b = append(b, 10) // add new line at the end

		if _, err := f.Write(b); err != nil {
			log.Printf("write failure: %s", err)
		}
		n++
	}
	return
}

func Load[E Entities](registry map[NumericID]E, fpath string) (n int) {
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

		switch any(e).(type) {
		case Balance:
			id := NumericID(reflect.ValueOf(e).FieldByName("Account").Int())
			registry[id] = e
			n++
		case Account, Transaction:
			id := NumericID(reflect.ValueOf(e).FieldByName("ID").Int())
			isDeleted := reflect.ValueOf(e).FieldByName("Deleted").Bool()
			if isDeleted {
				delete(registry, id)
			} else {
				registry[id] = e
				n++
			}
		}
	}
	return
}

// Remove all entites from journal marked for deletion.
func CleanUp[E Entities](fpath string) int {
	// todo delete from journal lines with Deleted: true
	return 0
}
