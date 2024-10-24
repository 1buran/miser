package miser

import (
	"sync"
	"time"
)

const Million = 1_000_000

const (
	Uncleared = iota // recorded but not yet reconciled; needs review
	Pending          // tentatively reconciled (if needed, eg during a big reconciliation)
	Cleared          // complete, reconciled as far as possible, and considered correct
)

type Transaction struct {
	ID, Source, Dest ID
	Time             time.Time
	Text             EncryptedString
	Value            int64 // in millionths
	State            int   // one of: Uncleared, Pending, Cleared
	Deleted          bool
}

func (t *Transaction) IsInitial() bool { return t.Source == t.Dest }

type TransactionRegistry struct {
	items  []Transaction
	queued []Transaction

	sync.RWMutex
}

func (tr *TransactionRegistry) List() (transactions []Transaction) {
	tr.RLock()
	defer tr.RUnlock()

	m := make(map[ID]int)

	for _, transa := range tr.items {
		n, oldVersion := m[transa.ID]
		if oldVersion {
			transactions[n] = transa
		} else {
			m[transa.ID] = len(transactions)
			transactions = append(transactions, transa)
		}
	}
	return transactions
}

func (tr *TransactionRegistry) Add(t Transaction) int {
	tr.Lock()
	defer tr.Unlock()
	tr.items = append(tr.items, t)
	return 1
}

func (tr *TransactionRegistry) AddQueued(t Transaction) {
	tr.Lock()
	defer tr.Unlock()
	tr.queued = append(tr.queued, t)
}

func (tr *TransactionRegistry) SyncQueued() []Transaction {
	tr.RLock()
	defer tr.RUnlock()
	return tr.queued
}

// Find last transaction.
func (tr *TransactionRegistry) Last(accID ID) *Transaction {
	var transa *Transaction
	var mostRecent time.Time

	for _, t := range tr.items {
		if (t.Source == accID || t.Dest == accID) && t.Time.After(mostRecent) {
			mostRecent = t.Time
			transa = &t
		}
	}
	return transa
}

// Find a transaction of account before given time.
func (tr *TransactionRegistry) FirstBefore(accID ID, trTime time.Time) *Transaction {
	tr.RLock()
	defer tr.RUnlock()

	for i := len(tr.items) - 1; i >= 0; i-- {
		t := tr.items[i]
		if (t.Source == accID || t.Dest == accID) && t.Time.Before(trTime) {
			return &t
		}
	}
	return nil
}

// Find all transactions of account after given time.
func (tr *TransactionRegistry) AllAfter(accID ID, trTime time.Time) (trs []Transaction) {
	tr.RLock()
	defer tr.RUnlock()

	for _, t := range tr.items {
		if (t.Source == accID || t.Dest == accID) && t.Time.After(trTime) {
			trs = append(trs, t)
		}
	}
	return
}

// Delete all transaction of given account (useful in case of account deletion).
// func DeleteAllAccountTransactions(accID ID) {
// 	for _, tr := range Transactions.items {
// 		if tr.Dest == accID || tr.Source == accID {
// 			DeleteTransaction(&tr)
// 		}
// 	}
// }

func CreateTransactionRegistry() *TransactionRegistry { return &TransactionRegistry{} }
func (tr *TransactionRegistry) Load() (int, error)    { return Load(tr, TRANSACTIONS_FILE) }
func (tr *TransactionRegistry) Save() (int, error)    { return Save(tr, TRANSACTIONS_FILE) }
