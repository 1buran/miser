package main

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
	_, _, symbol := Currency(string(Accounts[t.Source].Cur))
	return fmt.Sprintf("%c %.2f", symbol, float64(t.Value)/Million)
}

var Transactions map[NumericID]Transaction = make(map[NumericID]Transaction)
var syncTransactions map[NumericID]struct{} = make(map[NumericID]struct{})

func UpdateTransaction(tr *Transaction) { syncTransactions[tr.ID] = struct{}{} }

func DeleteTransaction(tr *Transaction) {
	tr.Deleted = true
	syncTransactions[tr.ID] = struct{}{}
}

func DeleteAllAccountTransactions(accID NumericID) {
	for _, tr := range Transactions {
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

	srcAcc, found := Accounts[src]
	if !found {
		return nil, errors.New("src account not found")
	}

	dstAcc, found := Accounts[dst]
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
	Transactions[tr.ID] = tr
	syncTransactions[tr.ID] = struct{}{}
	return &tr, nil
}

func LoadTransactions() int { return Load(Transactions, TRANSACTIONS_FILE) }
func SyncTransactions() int { return Save(syncTransactions, Transactions, TRANSACTIONS_FILE) }
