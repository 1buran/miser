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

	sync.RWMutex
}

func (ar *AccountRegistry) List() map[ID]Account {
	ar.RLock()
	defer ar.RUnlock()
	accounts := make(map[ID]Account)
	for _, acc := range ar.Items { // the last readed is the most actual version
		accounts[acc.ID] = acc
	}
	return accounts
}

func (ar *AccountRegistry) Get(accID ID) *Account {
	ar.RLock()
	defer ar.RUnlock()

	for i := len(ar.Items) - 1; i >= 0; i-- {
		item := ar.Items[i]
		if item.ID == accID {
			return &item
		}
	}
	return nil
}

func (ar *AccountRegistry) Add(a Account) int {
	ar.Lock()
	defer ar.Unlock()
	ar.Items = append(ar.Items, a)
	return 1

}

func (ar *AccountRegistry) AddQueued(a Account) {
	ar.Lock()
	defer ar.Unlock()
	ar.Queued = append(ar.Queued, a)
}

func (ar *AccountRegistry) SyncQueued() []Account {
	ar.RLock()
	defer ar.RUnlock()
	return ar.Queued
}

func CreateAccount(n, t, d, c string, initBalance float64) (*Account, error) {
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

	v := int64(initBalance * Million)
	tr := CreateInitialTransaction(acc.ID, v)
	b := CreateBalance(acc.ID, tr.ID, tr.Time, v)

	tag := Tags.GetByName(Initial)
	if tag == nil {
		tag = Tags.Create(Initial)
	}

	TagsMap.Create(tag.ID, tr.ID, TransactionTag)
	TagsMap.Create(tag.ID, b.ID, BalanceTag)

	return &acc, nil
}

var Accounts = AccountRegistry{}

func LoadAccounts() (int, error) { return Load(&Accounts, ACCOUNTS_FILE) }
func SyncAccounts() (int, error) { return Save(&Accounts, ACCOUNTS_FILE) }
