package miser

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Ledger struct {
	tr *TransactionRegistry
	ar *AccountRegistry
	br *BalanceRegistry
	cr *CurrencyRegistry
	tg *TagRegistry
	tm *TagMapRegistry
}

func CreateLedger(ar *AccountRegistry, br *BalanceRegistry, tr *TransactionRegistry, cr *CurrencyRegistry, tg *TagRegistry, tm *TagMapRegistry) *Ledger {
	return &Ledger{ar: ar, tr: tr, br: br, cr: cr, tg: tg, tm: tm}
}

// Save all queued data, sync it to disk.
func (l *Ledger) Save() {
	l.tr.Save()
	l.br.Save()
	l.ar.Save()
	l.tg.Save()
	l.tm.Save()
}

func (l *Ledger) CreateInitialTransaction(accID ID, openedAt time.Time, v int64) *Transaction {
	transa := Transaction{
		ID: CreateID(), Source: accID, Dest: accID, Time: openedAt,
		Value: v, Text: "Initial balance"}
	l.tr.Add(transa)
	l.tr.AddQueued(transa)

	return &transa
}

func (l *Ledger) CreateTransaction(src, dst ID, t time.Time, v float64, txt string) (*Transaction, error) {
	if v <= 0 {
		return nil, errors.New("transaction value should be greater zero")
	}

	if t.IsZero() {
		return nil, errors.New("zero time of transaction is not allowed")
	}

	srcAcc := l.ar.Get(src)
	if srcAcc == nil {
		return nil, errors.New("src account not found")
	}

	dstAcc := l.ar.Get(dst)
	if dstAcc == nil {
		return nil, errors.New("dst account not found")
	}

	if t.Before(srcAcc.OpenedAt) || t.Before(dstAcc.OpenedAt) {
		return nil, errors.New("transaction cannot be before the account is opened")
	}

	value := int64(v * Million)

	b := l.AccountBalance(src)
	if b.Value < value {
		return nil, errors.New("you cannot trasfer more money than you have")
	}

	if srcAcc.Type == dstAcc.Type {
		return nil, errors.New("cannot be transferred to same type of account")
	}

	transa := Transaction{
		ID:     CreateID(),
		Source: src,
		Dest:   dst,
		Time:   t,
		Value:  value,
		Text:   EncryptedString(txt),
	}
	l.tr.Add(transa)
	l.tr.AddQueued(transa)

	if err := l.UpdateBalance(src, transa.ID, string(srcAcc.Type), Credit, t, value); err != nil {
		return nil, err
	}
	if err := l.UpdateBalance(dst, transa.ID, string(dstAcc.Type), Debit, t, value); err != nil {
		return nil, err
	}

	return &transa, nil
}

func (l *Ledger) CreateAccount(n, t, d, c string, openedAt time.Time, initBalance float64) (*Account, error) {
	n = strings.TrimSpace(n)
	if n == "" {
		return nil, errors.New("name of account is blank")
	}

	if t != Asset && t != Liability && t != Equity && t != Income && t != Expense {
		return nil, fmt.Errorf("wrong type of account: %s", t)
	}

	if cur := l.cr.Get(c); cur == nil {
		return nil, fmt.Errorf("currency %q is not supproted", c)
	}

	acc := Account{
		ID:       CreateID(),
		Name:     EncryptedString(n),
		Type:     EncryptedString(t),
		Desc:     EncryptedString(d),
		Cur:      EncryptedString(c),
		OpenedAt: openedAt,
	}
	// add new account to registry and sync queue
	l.ar.Add(acc)
	l.ar.AddQueued(acc)

	// create initial transaction
	v := int64(initBalance * Million)
	transa := l.CreateInitialTransaction(acc.ID, openedAt, v)
	l.CreateBalance(acc.ID, transa.ID, v)

	// tag transaction as initial
	tag := l.tg.GetByName(Initial)
	if tag == nil {
		tag = l.tg.Create(Initial)
	}
	l.tm.Create(tag.ID, transa.ID)

	return &acc, nil
}

func (l *Ledger) CreateBalance(accID, trID ID, value int64) *Balance {
	b := Balance{Account: accID, Transaction: trID, Value: value}
	l.br.Add(b)
	l.br.AddQueued(b)
	return &b
}

// Credit - source, Debit - destination
func (l *Ledger) UpdateBalance(accID, trID ID, accType string, operType int, trTime time.Time, value int64) error {
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

	t := l.tr.FirstBefore(accID, trTime)
	if t == nil {
		return fmt.Errorf("transaction not found, before %s, account ID: %s", trTime, accID)
	}

	b := l.br.TransactionBalance(accID, t.ID)
	if b == nil {
		return fmt.Errorf("balance not found, transaction ID: %s, account ID: %s", t.ID, accID)
	}

	l.CreateBalance(accID, trID, b.Value+value)

	// rebalance in case if the current transaction was in the middle of history:
	//   t b      t b
	//   0 200    0 200
	// -50 150  -50 150
	// -20 130   -5 145  <-- a new transaction in the middle of history
	// -13 117  -20 125
	//          -13 112
	// fix = -5 (value of middle transaction)
	for _, transa := range l.tr.AllAfter(accID, trTime) {
		oldBalance := l.br.TransactionBalance(accID, transa.ID)
		if oldBalance != nil {
			l.CreateBalance(accID, transa.ID, oldBalance.Value+value)
		}
	}

	return nil
}

func (l *Ledger) AmountTransaction(t *Transaction) string {
	acc := l.ar.Get(t.Source)
	c := l.cr.Get(string(acc.Cur))
	if c != nil {
		return fmt.Sprintf("%s %.2f", c.Sign, float64(t.Value)/Million)
	}
	return fmt.Sprintf("%.2f", float64(t.Value)/Million)

}

// Account balance: balance at time of last transaction.
func (l *Ledger) AccountBalance(accID ID) *Balance {
	lastTransa := l.tr.Last(accID)
	if lastTransa != nil {
		b := l.br.TransactionBalance(accID, lastTransa.ID)
		if b != nil {
			return b
		}
	}
	return nil
}

// Account amount.
func (l *Ledger) AccountAmount(accID ID) float64 {
	b := l.AccountBalance(accID)
	if b == nil {
		return 0
	}
	return float64(b.Value) / Million
}
