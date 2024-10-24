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
	items  map[ID]Account
	queued map[ID]Account

	sync.RWMutex
}

func (ar *AccountRegistry) List() map[ID]Account {
	ar.RLock()
	defer ar.RUnlock()
	return ar.items
}

func (ar *AccountRegistry) Get(accID ID) *Account {
	ar.RLock()
	defer ar.RUnlock()
	acc, ok := ar.items[accID]
	if ok {
		return &acc
	}
	return nil
}

func (ar *AccountRegistry) Add(a Account) int {
	ar.Lock()
	defer ar.Unlock()
	ar.items[a.ID] = a
	return 1

}

func (ar *AccountRegistry) AddQueued(a Account) {
	ar.Lock()
	defer ar.Unlock()
	ar.queued[a.ID] = a
}

func (ar *AccountRegistry) SyncQueued() (changes []Account) {
	ar.RLock()
	defer ar.RUnlock()
	for _, acc := range ar.items {
		changes = append(changes, acc)
	}
	return
}

func CreateAccountRegistry() *AccountRegistry {
	return &AccountRegistry{items: make(map[ID]Account), queued: make(map[ID]Account)}
}

func (ar *AccountRegistry) Load() (int, error) { return Load(ar, ACCOUNTS_FILE) }
func (ar *AccountRegistry) Save() (int, error) { return Save(ar, ACCOUNTS_FILE) }
