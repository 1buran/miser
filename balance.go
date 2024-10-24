package miser

import (
	"fmt"
	"sync"
)

const (
	Credit = iota
	Debit
)

// Balance is value object, it is immutable,
// do not try to change it, do create another one instead,
// the last version will be used (see Add method of BalanceRegistry).
type Balance struct {
	Account, Transaction ID // in fact the id of balance item is transaction id
	Value                int64
}

func (b Balance) ID() string      { return fmt.Sprintf("%s-%s", b.Account, b.Transaction) }
func (b Balance) Amount() float64 { return float64(b.Value) / Million }

type BalanceRegistry struct {
	items  []Balance          // order matters, items loaded from disk, last their versions
	queued map[string]Balance // queue of items for sync to disk

	idx map[string]int

	sync.RWMutex
}

func (br *BalanceRegistry) List() []Balance {
	br.RLock()
	defer br.RUnlock()
	return br.items
}

// Find a balance of given account transaction.
// Keep in mind: every transaction creates two different balances
// for source and destination accounts.
func (br *BalanceRegistry) TransactionBalance(accID, trID ID) *Balance {
	br.RLock()
	defer br.RUnlock()

	key := fmt.Sprintf("%s-%s", accID, trID)
	i, ok := br.idx[key]
	if ok {
		return &br.items[i]
	}
	return nil
}

func (br *BalanceRegistry) Add(b Balance) int {
	br.Lock()
	defer br.Unlock()

	// check the index of item and update it in place
	i, ok := br.idx[b.ID()]
	if ok {
		br.items[i] = b
		return 1
	}

	br.idx[b.ID()] = len(br.items)
	br.items = append(br.items, b)
	return 1
}

func (br *BalanceRegistry) AddQueued(b Balance) {
	br.Lock()
	defer br.Unlock()
	br.queued[b.ID()] = b
}

func (br *BalanceRegistry) SyncQueued() (changes []Balance) {
	br.RLock()
	defer br.RUnlock()
	for _, v := range br.queued {
		changes = append(changes, v)
	}
	return
}

func CreateBalanceRegistry() *BalanceRegistry {
	return &BalanceRegistry{idx: make(map[string]int), queued: make(map[string]Balance)}
}

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
