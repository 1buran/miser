package main

import (
	"time"
)

const (
	Credit = iota
	Debit
)

type Balance struct {
	Account      NumericID
	Value        int64
	ReconciledAt time.Time
}

func (b *Balance) isReconciled() bool { return !b.ReconciledAt.IsZero() }
func (b *Balance) setReconciled() {
	b.ReconciledAt = time.Now()
	syncBalances[b.Account] = struct{}{}
}

var Balances map[NumericID]Balance = make(map[NumericID]Balance)
var syncBalances map[NumericID]struct{} = make(map[NumericID]struct{})

func LoadBalances() int { return Load(Balances, BALANCE_FILE) }
func SaveBalances() int { return Save(syncBalances, Balances, BALANCE_FILE) }

// Account Type  | Effect on Account Balance
// ------------------------------------------
// --------------|    Debit     |   Credit
// --------------|---------------------------
// Assets        |              |
//               | Increase     |  Decrease
// Expenses      |              |
// -------------------------------------------
// Liabilities   |              |
// Equity        | Decrease     |  Increase
// Income        |              |
// -------------------------------------------
//

// Credit - source, Debit - destination
func UpdateBalance(accID NumericID, accType string, operType int, value int64) {
	switch operType {
	case Credit:
		if accType == Asset || accType == Expense {
			value = -value
		}
	case Debit:
		if accType == Liability || accType == Equity || accType == Income {
			value = -value
		}
	}

	b, found := Balances[accID]
	if !found {
		b = Balance{Account: accID, Value: value}
		Balances[accID] = b
		return
	}
	b.Value += value
}

func RefreshBalances() (n int) {
	for _, tr := range Transactions {
		if Balances[tr.Source].ReconciledAt.Before(tr.Time) {
			UpdateBalance(tr.Source, string(Accounts[tr.Source].Type), Credit, tr.Value)
			n++
		}

		if Balances[tr.Dest].ReconciledAt.Before(tr.Time) {
			UpdateBalance(tr.Dest, string(Accounts[tr.Dest].Type), Debit, tr.Value)
			n++
		}
	}
	return
}

// The rearranged accounting equation:
// Assets + Expenses = Liabilities + Equity + Income
func CheckBalance() int64 {
	var as, li, eq, in, ex int64
	for accID, bl := range Balances {
		switch Accounts[accID].Type {
		case Asset:
			as += bl.Value
		case Liability:
			li += bl.Value
		case Equity:
			eq += bl.Value
		case Income:
			in += bl.Value
		case Expense:
			ex += bl.Value
		}
	}
	return as + ex - li - eq - in
}
