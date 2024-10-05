package miser

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	Asset     = "Asset"
	Liability = "Liability"
	Equity    = "Equity"
	Income    = "Income"
	Expense   = "Expense"
)

type Account struct {
	ID NumericID `json:"id"`

	Name EncryptedString `json:"Name"`
	Type EncryptedString `json:"Type"`
	Desc EncryptedString `json:"Desc"`
	Cur  EncryptedString `json:"Cur"`

	OpenedAt time.Time `json:"OpenedAt"`
	ClosedAt time.Time `json:"ClosedAt"`

	Deleted bool `json:"Deleted"`
}

func (a *Account) isClosed() bool { return !a.ClosedAt.IsZero() }

var Accounts map[NumericID]Account = make(map[NumericID]Account)
var syncAccounts map[NumericID]struct{} = make(map[NumericID]struct{})

func UpdateAccount(acc *Account) { syncAccounts[acc.ID] = struct{}{} }

func DeleteAccount(acc *Account) {
	acc.Deleted = true
	syncAccounts[acc.ID] = struct{}{}

	DeleteAllAccountTransactions(acc.ID)
}

func CreateAccount(n, t, d, c string) (*Account, error) {
	n = strings.TrimSpace(n)
	if n == "" {
		return nil, errors.New("name of account is blank")
	}

	if t != Asset && t != Liability && t != Equity && t != Income && t != Expense {
		return nil, fmt.Errorf("wrong type of account: %s", t)
	}

	if supported, _, _ := Currency(c); !supported {
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
	Accounts[acc.ID] = acc
	syncAccounts[acc.ID] = struct{}{}

	return &acc, nil
}

func LoadAccounts() int { return Load(Accounts, ACCOUNTS_FILE) }
func SyncAccounts() int { return Save(syncAccounts, Accounts, ACCOUNTS_FILE) }
