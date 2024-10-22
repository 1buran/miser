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
}

func CreateLedger(ar *AccountRegistry, br *BalanceRegistry, tr *TransactionRegistry, cr *CurrencyRegistry) *Ledger {
	return &Ledger{ar: ar, tr: tr, br: br, cr: cr}
}

func (l *Ledger) CreateInitialTransaction(accID ID, v int64) *Transaction {
	transa := Transaction{
		ID: CreateID(), Source: accID, Dest: accID, Time: time.Now(),
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

	b := l.br.AccountBalance(src)
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

	l.UpdateBalance(src, transa.ID, string(srcAcc.Type), Credit, transa.Time, value)
	l.UpdateBalance(dst, transa.ID, string(dstAcc.Type), Debit, transa.Time, value)

	return &transa, nil
}

func (l *Ledger) CreateAccount(n, t, d, c string, initBalance float64) (*Account, error) {
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
		OpenedAt: time.Now(),
	}
	l.ar.Add(acc)
	l.ar.AddQueued(acc)

	v := int64(initBalance * Million)
	transa := l.CreateInitialTransaction(acc.ID, v)
	b := l.CreateBalance(acc.ID, transa.ID, transa.Time, v)

	tag := Tags.GetByName(Initial)
	if tag == nil {
		tag = Tags.Create(Initial)
	}

	TagsMap.Create(tag.ID, transa.ID, TransactionTag)
	TagsMap.Create(tag.ID, b.ID, BalanceTag)

	return &acc, nil

}

func (l *Ledger) CreateBalance(accID, trID ID, t time.Time, value int64) *Balance {
	b := Balance{ID: CreateID(), Account: accID, Transaction: trID, Value: value, Time: time.Now()}
	l.br.Add(b)
	l.br.AddQueued(b)
	return &b
}

// Credit - source, Debit - destination
func (l *Ledger) UpdateBalance(accID, trID ID, accType string, operType int, trTime time.Time, value int64) {
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

	b := l.br.AccountBalanceBefore(accID, trTime)
	if b == nil {
		panic("tried update nil balance, account ID:" + accID)
	}

	value += b.Value

	l.CreateBalance(accID, trID, trTime, value)

	b2 := l.br.AccountBalanceAfter(accID, trTime)
	if b2 != nil {
		// new delta between balance changes
		fix := value - b2.Value

		for _, change := range l.br.rebalanceAccountBalanceAfter(accID, trTime, fix) {
			l.br.Add(*change)
			l.br.AddQueued(*change)
		}
	}
}

func (l *Ledger) AmountTransaction(t *Transaction) string {
	acc := l.ar.Get(t.Source)
	c := l.cr.Get(string(acc.Cur))
	if c != nil {
		return fmt.Sprintf("%s %.2f", c.Sign, float64(t.Value)/Million)
	}
	return fmt.Sprintf("%.2f", float64(t.Value)/Million)
}
