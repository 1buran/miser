package miser

import (
	"errors"
	"fmt"
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
func (t *Transaction) Amount() string {
	acc := Accounts.Get(t.Source)
	_, _, symbol := Currency(string(acc.Cur))
	return fmt.Sprintf("%c %.2f", symbol, float64(t.Value)/Million)
}

func UpdateTransaction(t *Transaction) { Transactions.AddQueued(*t) }

func DeleteTransaction(t *Transaction) {
	t.Deleted = true
	Transactions.AddQueued(*t)
}

type TransactionRegistry struct {
	Items  []Transaction
	Queued []Transaction

	sync.RWMutex
}

func (tr *TransactionRegistry) List() (transactions []Transaction) {
	tr.RLock()
	defer tr.RUnlock()

	m := make(map[ID]int)

	for _, transa := range tr.Items {
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
	tr.Items = append(tr.Items, t)
	return 1
}

func (tr *TransactionRegistry) AddQueued(t Transaction) {
	tr.Lock()
	defer tr.Unlock()
	tr.Queued = append(tr.Queued, t)
}

func (tr *TransactionRegistry) SyncQueued() []Transaction {
	tr.RLock()
	defer tr.RUnlock()
	return tr.Queued
}

// Delete all transaction of given account (useful in case of account deletion).
func DeleteAllAccountTransactions(accID ID) {
	for _, tr := range Transactions.Items {
		if tr.Dest == accID || tr.Source == accID {
			DeleteTransaction(&tr)
		}
	}
}

func CreateTransation(src, dst ID, t time.Time, v float64, txt string) (*Transaction, error) {
	if v <= 0 {
		return nil, errors.New("transaction value should be greater zero")
	}

	if t.IsZero() {
		return nil, errors.New("zero time of transaction is not allowed")
	}

	srcAcc := Accounts.Get(src)
	if srcAcc == nil {
		return nil, errors.New("src account not found")
	}

	dstAcc := Accounts.Get(dst)
	if dstAcc == nil {
		return nil, errors.New("dst account not found")
	}

	if t.Before(srcAcc.OpenedAt) || t.Before(dstAcc.OpenedAt) {
		return nil, errors.New("transaction cannot be before the account is opened")
	}

	value := int64(v * Million)

	b := Balances.AccountBalance(src)
	if b.Value < value {
		return nil, errors.New("you cannot trasfer more money than you have")
	}

	if srcAcc.Type == dstAcc.Type {
		return nil, errors.New("cannot be transferred to same type of account")
	}

	tr := Transaction{
		ID:     CreateID(),
		Source: src,
		Dest:   dst,
		Time:   t,
		Value:  value,
		Text:   EncryptedString(txt),
	}
	Transactions.Add(tr)
	Transactions.AddQueued(tr)

	UpdateBalance(src, tr.ID, string(srcAcc.Type), Credit, tr.Time, value)
	UpdateBalance(dst, tr.ID, string(dstAcc.Type), Debit, tr.Time, value)

	return &tr, nil
}

var Transactions = TransactionRegistry{}

func LoadTransactions() (int, error) { return Load(&Transactions, TRANSACTIONS_FILE) }
func SyncTransactions() (int, error) { return Save(&Transactions, TRANSACTIONS_FILE) }

func CreateInitialTransaction(accID ID, v int64) *Transaction {
	tr := Transaction{
		ID: CreateID(), Source: accID, Dest: accID, Time: time.Now(),
		Value: v, Text: "Initial balance"}
	Transactions.Add(tr)
	Transactions.AddQueued(tr)

	return &tr
}
