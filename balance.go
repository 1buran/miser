package miser

import (
	"sync"
	"time"
)

const (
	Credit = iota
	Debit
)

type Balance struct {
	ID, Account, Transaction ID
	Value                    int64
	Time                     time.Time
	Deleted                  bool
}

type BalanceRegistry struct {
	items  []Balance // all items loaded from disk
	queued []Balance // queue of items for sync to disk

	sync.RWMutex
}

func (br *BalanceRegistry) List() (balances []Balance) {
	br.RLock()
	defer br.RUnlock()

	m := make(map[ID]int)

	for _, balance := range br.items {
		if balance.Deleted {
			continue
		}
		n, oldVer := m[balance.ID]
		if oldVer {
			balances[n] = balance
		} else {
			m[balance.ID] = len(balances)
			balances = append(balances, balance)
		}
	}
	return
}

// Find a balance of given account transaction.
func (br *BalanceRegistry) TransactionBalance(accID, trID ID) *Balance {
	br.RLock()
	defer br.RUnlock()
	for i := len(br.items) - 1; i >= 0; i-- {
		item := br.items[i]
		if !item.Deleted && item.Account == accID && item.Transaction == trID {
			return &item
		}
	}
	return nil
}

// Find a current(last) balance of account.
func (br *BalanceRegistry) AccountBalance(accID ID) *Balance {
	br.RLock()
	defer br.RUnlock()
	for i := len(br.items) - 1; i >= 0; i-- {
		if !br.items[i].Deleted && br.items[i].Account == accID {
			return &br.items[i]
		}
	}
	return nil
}

// Find a balance of account before given time.
func (br *BalanceRegistry) AccountBalanceBefore(accID ID, t time.Time) *Balance {
	br.RLock()
	defer br.RUnlock()
	for i := len(br.items) - 1; i >= 0; i-- {
		if !br.items[i].Deleted && br.items[i].Account == accID && br.items[i].Time.Before(t) {
			return &br.items[i]
		}
	}
	return nil
}

// Find a balance of account after given time.
func (br *BalanceRegistry) AccountBalanceAfter(accID ID, t time.Time) *Balance {
	br.RLock()
	defer br.RUnlock()
	for i := 0; i < len(br.items); i++ {
		if !br.items[i].Deleted && br.items[i].Account == accID && br.items[i].Time.After(t) {
			return &br.items[i]
		}
	}
	return nil
}

// Rebalance account balance after changes.
func (br *BalanceRegistry) rebalanceAccountBalanceAfter(accID ID, t time.Time, fix int64) (changes []*Balance) {
	br.RLock()
	defer br.RUnlock()
	for i := 0; i < len(br.items); i++ {
		if !br.items[i].Deleted && br.items[i].Account == accID && br.items[i].Time.After(t) {
			b := br.items[i]
			b.Deleted = true // mark old one as deleted
			changes = append(
				changes, &b, &Balance{ // create another one with fixed balance
					ID: CreateID(), Account: b.Account, Transaction: b.Transaction,
					Time: b.Time, Value: b.Value + fix})
		}
	}
	return
}

func (br *BalanceRegistry) AccountValue(accID ID) float64 {
	b := br.AccountBalance(accID)
	if b == nil {
		return 0
	}
	return float64(b.Value) / Million
}

// Update balance in registry.
// todo considering change registry items to map of pointers for able update items directly
func (br *BalanceRegistry) Update(b Balance) {
	br.Lock()
	defer br.Unlock()
	for i := len(br.items) - 1; i >= 0; i-- {
		item := br.items[i]
		if item.Deleted {
			continue
		}
		if item.ID == b.ID {
			item = b
			break
		}
	}
	br.AddQueued(b)
}

func (br *BalanceRegistry) Add(b Balance) int {
	br.Lock()
	defer br.Unlock()
	br.items = append(br.items, b)
	return 1
}

func (br *BalanceRegistry) AddQueued(b Balance) {
	br.Lock()
	defer br.Unlock()
	br.queued = append(br.queued, b)
}

func (br *BalanceRegistry) SyncQueued() []Balance {
	br.RLock()
	defer br.RUnlock()
	return br.queued
}

func CreateBalanceRegistry() *BalanceRegistry  { return &BalanceRegistry{} }
func (br *BalanceRegistry) Load() (int, error) { return Load(br, BALANCE_FILE) }
func (br *BalanceRegistry) Save() (int, error) { return Save(br, BALANCE_FILE) }

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
