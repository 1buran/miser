package miser

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const Million = 1_000_000

const (
	Uncleared = iota // recorded but not yet reconciled; needs review
	Pending          // tentatively reconciled (if needed, eg during a big reconciliation)
	Cleared          // complete, reconciled as far as possible, and considered correct
)

type Transaction struct {
	ID, Source, Dest NumericID
	Time             time.Time
	Text             EncryptedString
	Value            int64 // in millionths
	State            int   // one of: Uncleared, Pending, Cleared
	Deleted          bool
}

func (t *Transaction) Amount() string {
	acc, _ := Accounts.Get(t.Source)
	_, _, symbol := Currency(string(acc.Cur))
	return fmt.Sprintf("%c %.2f", symbol, float64(t.Value)/Million)
}

func UpdateTransaction(t *Transaction) { Transactions.AddQueued(*t) }

func DeleteTransaction(t *Transaction) {
	t.Deleted = true
	Transactions.AddQueued(*t)
}

type TransactionRegistry struct {
	Items  map[NumericID]Transaction
	Queued map[NumericID]struct{}
}

func (tr TransactionRegistry) Add(t Transaction) int {
	if t.Deleted {
		delete(tr.Items, t.ID)
	} else {
		tr.Items[t.ID] = t
		return 1
	}
	return 0
}

func (tr TransactionRegistry) AddQueued(t Transaction) {
	tr.Queued[t.ID] = struct{}{}
}

func (tr TransactionRegistry) SyncQueued() []Transaction {
	var items []Transaction
	for id := range tr.Queued {
		items = append(items, tr.Items[id])
	}
	return items
}

func DeleteAllAccountTransactions(accID NumericID) {
	for _, tr := range Transactions.Items {
		if tr.Dest == accID || tr.Source == accID {
			DeleteTransaction(&tr)
		}
	}
}

func CreateTransation(src, dst NumericID, t time.Time, v string, txt string) (*Transaction, error) {
	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, err
	}

	if val <= 0 {
		return nil, errors.New("transaction value should be greater zero")
	}

	if t.IsZero() {
		return nil, errors.New("zero time of transaction is not allowed")
	}

	srcAcc, found := Accounts.Get(src)
	if !found {
		return nil, errors.New("src account not found")
	}

	dstAcc, found := Accounts.Get(dst)
	if !found {
		return nil, errors.New("dst account not found")
	}

	value := int64(val * Million)

	if srcAcc.Type == dstAcc.Type {
		return nil, errors.New("cannot be transferred to same type of account")
	}

	UpdateBalance(src, string(srcAcc.Type), Credit, value)
	UpdateBalance(dst, string(dstAcc.Type), Debit, value)

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

	return &tr, nil
}

var Transactions = TransactionRegistry{
	Items:  make(map[NumericID]Transaction),
	Queued: make(map[NumericID]struct{}),
}

func LoadTransactions() int { return Load(Transactions, TRANSACTIONS_FILE) }
func SyncTransactions() int { return Save(Transactions, TRANSACTIONS_FILE) }
