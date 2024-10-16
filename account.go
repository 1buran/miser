package miser

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	Asset     = "Asset"
	Liability = "Liability"
	Equity    = "Equity"
	Income    = "Income"
	Expense   = "Expense"
)

var accMu, crossMu sync.Mutex

type Account struct {
	ID                    ID
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
	Items  []Account
	Queued []Account
}

func (ar AccountRegistry) Get(accID ID) *Account {
	accMu.Lock()
	defer accMu.Unlock()
	for _, item := range ar.Items {
		if item.ID == accID {
			return &item
		}
	}
	return nil
}

func (ar *AccountRegistry) Add(a Account) int {
	accMu.Lock()
	defer accMu.Unlock()
	ar.Items = append(ar.Items, a)
	return 1

}

func (ar *AccountRegistry) AddQueued(a Account) {
	accMu.Lock()
	defer accMu.Unlock()
	ar.Queued = append(ar.Queued, a)
}

func (ar AccountRegistry) SyncQueued() []Account {
	accMu.Lock()
	defer accMu.Unlock()
	return ar.Queued
}

func CreateAccount(n, t, d, c string, initBalance float64) (*Account, error) {
	crossMu.Lock()
	defer crossMu.Unlock()

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

	CreateBalance(acc.ID, int64(initBalance*Million))

	return &acc, nil
}

var Accounts = AccountRegistry{}

func LoadAccounts() (int, error) { return Load(&Accounts, ACCOUNTS_FILE) }
func SyncAccounts() (int, error) { return Save(&Accounts, ACCOUNTS_FILE) }
