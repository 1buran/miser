package miser

import (
	"sync"
	"time"
)

const (
	Credit = iota
	Debit
)

var balMu sync.Mutex

type Balance struct {
	ID, Account, Transaction ID
	Value                    int64
	Time                     time.Time
}

type BalanceRegistry struct {
	Items  []Balance // all items loaded from disk
	Queued []Balance // queue of items for sync to disk
}

// Find a balance of given account transaction.
func (br BalanceRegistry) TransactionBalance(accID, trID ID) *Balance {
	balMu.Lock()
	defer balMu.Unlock()
	for i := len(br.Items) - 1; i >= 0; i-- {
		item := br.Items[i]
		if item.Account == accID && item.Transaction == trID {
			return &item
		}
	}
	return nil
}

// Find a current(last) balance of account.
func (br BalanceRegistry) AccountBalance(accID ID) *Balance {
	balMu.Lock()
	defer balMu.Unlock()
	for i := len(br.Items) - 1; i >= 0; i-- {
		if br.Items[i].Account == accID {
			return &br.Items[i]
		}
	}
	return nil
}

func (br BalanceRegistry) AccountValue(accID ID) float64 {
	b := br.AccountBalance(accID)
	if b == nil {
		return 0
	}
	return float64(b.Value) / Million
}

// Update balance in registry.
// todo considering change registry items to map of pointers for able update items directly
func (br BalanceRegistry) Update(b Balance) {
	balMu.Lock()
	defer balMu.Unlock()
	for i := len(br.Items) - 1; i >= 0; i-- {
		item := br.Items[i]
		if item.ID == b.ID {
			item = b
			break
		}
	}
	br.AddQueued(b)
}

func (br *BalanceRegistry) Add(b Balance) int {
	balMu.Lock()
	defer balMu.Unlock()
	br.Items = append(br.Items, b)
	return 1
}

func (br *BalanceRegistry) AddQueued(b Balance) {
	balMu.Lock()
	defer balMu.Unlock()
	br.Queued = append(br.Queued, b)
}

func (br BalanceRegistry) SyncQueued() []Balance {
	balMu.Lock()
	defer balMu.Unlock()
	return br.Queued
}

var Balances = BalanceRegistry{}

func LoadBalances() (int, error) { return Load(&Balances, BALANCE_FILE) }
func SaveBalances() (int, error) { return Save(&Balances, BALANCE_FILE) }

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
func UpdateBalance(accID ID, accType string, operType int, value int64) {
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

	b := Balances.AccountBalance(accID)
	if b == nil {
		panic("tried update nil balance, account ID:" + accID)
	}
	value += b.Value
	CreateBalance(accID, value)
}

// Create balance for an account
func CreateBalance(accID ID, value int64) {
	b := Balance{ID: CreateID(), Account: accID, Value: value, Time: time.Now()}
	Balances.Add(b)
	Balances.AddQueued(b)
}

// The rearranged accounting equation:
// Assets + Expenses = Liabilities + Equity + Income
// func CheckBalance() int64 {
// 	var as, li, eq, in, ex int64
// 	for accID, bl := range Balances.Items {
// 		acc, _ := Accounts.Get(accID)
// 		switch acc.Type {
// 		case Asset:
// 			as += bl.Value
// 		case Liability:
// 			li += bl.Value
// 		case Equity:
// 			eq += bl.Value
// 		case Income:
// 			in += bl.Value
// 		case Expense:
// 			ex += bl.Value
// 		}
// 	}
// 	return as + ex - li - eq - in
// }
