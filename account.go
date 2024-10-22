package miser

import (
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

type AccountRegistry struct {
	items  []Account
	queued []Account

	sync.RWMutex
}

func (ar *AccountRegistry) List() map[ID]Account {
	ar.RLock()
	defer ar.RUnlock()
	accounts := make(map[ID]Account)
	for _, acc := range ar.items { // the last readed is the most actual version
		accounts[acc.ID] = acc
	}
	return accounts
}

func (ar *AccountRegistry) Get(accID ID) *Account {
	ar.RLock()
	defer ar.RUnlock()

	for i := len(ar.items) - 1; i >= 0; i-- {
		item := ar.items[i]
		if item.ID == accID {
			return &item
		}
	}
	return nil
}

func (ar *AccountRegistry) Add(a Account) int {
	ar.Lock()
	defer ar.Unlock()
	ar.items = append(ar.items, a)
	return 1

}

func (ar *AccountRegistry) AddQueued(a Account) {
	ar.Lock()
	defer ar.Unlock()
	ar.queued = append(ar.queued, a)
}

func (ar *AccountRegistry) SyncQueued() []Account {
	ar.RLock()
	defer ar.RUnlock()
	return ar.queued
}

func CreateAccountRegistry() *AccountRegistry  { return &AccountRegistry{} }
func (ar *AccountRegistry) Load() (int, error) { return Load(ar, ACCOUNTS_FILE) }
func (ar *AccountRegistry) Save() (int, error) { return Save(ar, ACCOUNTS_FILE) }
