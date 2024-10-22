package miser

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
)

const (
	ACCOUNTS_FILE     = "miser.ar"
	TRANSACTIONS_FILE = "miser.tr"
	BALANCE_FILE      = "miser.br"
	TAGS_FILE         = "miser.tg"
	TAGS_MAPPING_FILE = "miser.tm"
)

type Entities interface {
	Account | Transaction | Balance | Tag | TagMap
}

type Registry[E Entities] interface {
	*AccountRegistry | *TransactionRegistry | *BalanceRegistry | *TagRegistry | *TagMapRegistry

	Add(e E) int
	SyncQueued() []E
}

func Save[E Entities, R Registry[E]](registry R, fpath string) (n int, err error) {
	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return n, err
	}

	defer func() {
		if e := f.Close(); e != nil {
			if err != nil {
				err = errors.Join(e, err)
			} else {
				err = e
			}
		}
	}()

	for _, item := range registry.SyncQueued() {
		b, err := json.Marshal(item)
		if err != nil {
			return n, err
		}

		b = append(b, 10) // add new line at the end

		if _, err := f.Write(b); err != nil {
			return n, err
		}
		n++
	}
	return n, err
}

func Load[E Entities, R Registry[E]](registry R, fpath string) (n int, err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return n, err
	}

	defer func() {
		if e := f.Close(); e != nil {
			if err != nil {
				err = errors.Join(e, err)
			} else {
				err = e
			}
		}
	}()

	r := bufio.NewReader(f)

	for {
		b, err := r.ReadBytes(10)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return n, err
			}
		}
		var e E
		if err := json.Unmarshal(b, &e); err != nil {
			return n, err
		}
		n += registry.Add(e)
	}
	return n, err
}

// Remove all entites from journal marked for deletion.
func CleanUp[E Entities](fpath string) int {
	// todo delete from journal lines with Deleted: true
	return 0
}
