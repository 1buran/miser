package miser

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
	Balances.AddQueued(*b)
}

type BalanceRegistry struct {
	Items  map[NumericID]Balance
	Queued map[NumericID]struct{}
}

func (br BalanceRegistry) Get(id NumericID) (Balance, bool) {
	bl, ok := br.Items[id]
	return bl, ok
}

func (br BalanceRegistry) Add(b Balance) int {
	br.Items[b.Account] = b
	return 1
}

func (br BalanceRegistry) AddQueued(b Balance) {
	br.Queued[b.Account] = struct{}{}
}

func (br BalanceRegistry) SyncQueued() []Balance {
	var items []Balance
	for id := range br.Queued {
		items = append(items, br.Items[id])
	}
	return items
}

var Balances = BalanceRegistry{
	Items:  make(map[NumericID]Balance),
	Queued: make(map[NumericID]struct{}),
}

func LoadBalances() (int, error) { return Load(Balances, BALANCE_FILE) }
func SaveBalances() (int, error) { return Save(Balances, BALANCE_FILE) }

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

	b, found := Balances.Get(accID)
	if !found {
		b = Balance{Account: accID, Value: value}
		Balances.Add(b)
		return
	}
	b.Value += value

	Balances.AddQueued(b)
}

// The rearranged accounting equation:
// Assets + Expenses = Liabilities + Equity + Income
func CheckBalance() int64 {
	var as, li, eq, in, ex int64
	for accID, bl := range Balances.Items {
		acc, _ := Accounts.Get(accID)
		switch acc.Type {
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
