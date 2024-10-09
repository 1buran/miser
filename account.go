package miser

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	Asset     = "Asset"
	Liability = "Liability"
	Equity    = "Equity"
	Income    = "Income"
	Expense   = "Expense"
)

type Account struct {
	ID                    NumericID
	Name, Type, Desc, Cur EncryptedString
	OpenedAt, ClosedAt    time.Time
	Deleted               bool
}

func (a *Account) isClosed() bool { return !a.ClosedAt.IsZero() }

func UpdateAccount(acc *Account) { Accounts.AddQueued(*acc) }

func DeleteAccount(acc *Account) {
	acc.Deleted = true
	Accounts.AddQueued(*acc)
	DeleteAllAccountTransactions(acc.ID)
}

type AccountRegistry struct {
	Items  map[NumericID]Account
	Queued map[NumericID]struct{}
}

func (ar AccountRegistry) Get(id NumericID) (Account, bool) {
	acc, ok := ar.Items[id]
	return acc, ok
}

func (ar AccountRegistry) Add(a Account) int {
	if a.Deleted {
		delete(ar.Items, a.ID)
	} else {
		ar.Items[a.ID] = a
		return 1
	}
	return 0
}
func (ar AccountRegistry) AddQueued(a Account) {
	ar.Queued[a.ID] = struct{}{}
}

func (ar AccountRegistry) SyncQueued() []Account {
	var items []Account
	for id := range ar.Queued {
		items = append(items, ar.Items[id])
	}
	return items
}

func CreateAccount(n, t, d, c string) (*Account, error) {
	n = strings.TrimSpace(n)
	if n == "" {
		return nil, errors.New("name of account is blank")
	}

	if t != Asset && t != Liability && t != Equity && t != Income && t != Expense {
		return nil, fmt.Errorf("wrong type of account: %s", t)
	}

	if supported, _, _ := Currency(c); !supported {
		return nil, fmt.Errorf("currency %q is not supproted", c)
	}

	acc := Account{
		ID:       CreateID(),
		Name:     EncryptedString(n),
		Type:     EncryptedString(t),
		Desc:     EncryptedString(d),
		Cur:      EncryptedString(c),
		OpenedAt: time.Now(),
	}
	Accounts.Add(acc)
	Accounts.AddQueued(acc)

	return &acc, nil
}

var Accounts = AccountRegistry{
	Items:  make(map[NumericID]Account),
	Queued: make(map[NumericID]struct{}),
}

func LoadAccounts() int { return Load(Accounts, ACCOUNTS_FILE) }
func SyncAccounts() int { return Save(Accounts, ACCOUNTS_FILE) }
