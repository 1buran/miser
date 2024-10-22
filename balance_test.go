package miser

import (
	"testing"
	"time"
)

func TestInitBalance(t *testing.T) {

	t.Parallel()

	// Create repositories:
	ar := CreateAccountRegistry()
	tr := CreateTransactionRegistry()
	br := CreateBalanceRegistry()

	// Create service:
	l := CreateLedger(ar, br, tr)

	t.Run("zero", func(t *testing.T) {
		acc, err := l.CreateAccount("Deposit", Asset, "deposit account", "USD", 0.00)
		if err != nil {
			t.Fatal(err)
		}

		b := l.br.AccountBalance(acc.ID)
		if b == nil {
			t.Fatal("balance was not created during account creation")
		}

		if b.Value != 0 {
			t.Errorf("expected 0, %d found", b.Value)
		}

		tmaps := TagsMap.GetByItemId(b.ID)
		if len(tmaps) != 1 {
			t.Fatalf("expected 1 tag, found: %d", len(tmaps))
		}

		expectedTag := Tags.GetByName(Initial)
		if tmaps[0].Tag != expectedTag.ID {
			t.Errorf("unexpected tag found: %#v", Tags.GetById(tmaps[0].Tag))
		}

	})

	t.Run("float", func(t *testing.T) {
		acc, err := l.CreateAccount("Deposit", Asset, "deposit account", "USD", 123.78)
		if err != nil {
			t.Fatal(err)
		}

		if amount := l.br.AccountValue(acc.ID); amount != 123.78 {
			t.Errorf("expected 123.78, got: %f", amount)
		}
	})
}

func TestChangeBalance(t *testing.T) {
	t.Parallel()

	// Create repositories:
	ar := CreateAccountRegistry()
	tr := CreateTransactionRegistry()
	br := CreateBalanceRegistry()

	// Create service:
	l := CreateLedger(ar, br, tr)

	t.Run("Expense", func(t *testing.T) {
		cash, err := l.CreateAccount("Cash", Asset, "wallet", "USD", 1555.12)
		if err != nil {
			t.Fatal(err)
		}

		market, err := l.CreateAccount("Market", Expense, "holiday market", "USD", 343.11)
		if err != nil {
			t.Fatal(err)
		}

		if openMarketBalance := l.br.AccountValue(market.ID); openMarketBalance != 343.11 {
			t.Errorf("expected 344.64, got: %.2f", openMarketBalance)
		}

		_, err = l.CreateTransaction(cash.ID, market.ID, time.Now(), 1.53, "1kg carrot")
		if err != nil {
			t.Fatal(err)
		}

		if amount := l.br.AccountValue(cash.ID); amount != 1553.59 {
			t.Errorf("expected 1553.59, got: %.2f", amount)
		}

		if updateMarketBalance := l.br.AccountValue(market.ID); updateMarketBalance != 344.64 {
			t.Errorf("expected 344.64, got: %.2f", updateMarketBalance)
		}

	})

	t.Run("Earnings", func(t *testing.T) {
		var buyers []Account
		for i := 0; i < 5; i++ {
			acc, _ := l.CreateAccount("Cash", Asset, "wallet", "USD", 1555.12)
			buyers = append(buyers, *acc)
		}

		market, err := l.CreateAccount("Market", Expense, "holiday market", "USD", 343.11)
		if err != nil {
			t.Fatal(err)
		}

		openMarketBalance := l.br.AccountBalance(market.ID)
		t.Logf("open market balance: %#v", openMarketBalance)
		for i := 0; i < 5; i++ { // buy 5 kg of carrot
			_, _ = l.CreateTransaction(buyers[i].ID, market.ID, time.Now(), 1.53, "1kg carrot")
		}

		closeMarketBalance := l.br.AccountBalance(market.ID)
		t.Logf("close market balance: %#v", closeMarketBalance)

		earnings := float64(closeMarketBalance.Value-openMarketBalance.Value) / Million
		if earnings != 7.65 {
			t.Logf("%#v", earnings)
			t.Errorf("expected earnings for sold 5kg of carrot: %+.2f, got: %+.2f", 5*1.53, earnings)
		}
	})

	t.Run("Spendings", func(t *testing.T) {
		var pubs []Account

		for i := 0; i < 5; i++ {
			acc, _ := l.CreateAccount("Shop", Asset, "LC Waikiki", "USD", 15.12)
			pubs = append(pubs, *acc)
		}

		cash, err := l.CreateAccount("Cash", Expense, "wallet", "USD", 343.11)
		if err != nil {
			t.Fatal(err)
		}

		openPartyBalance := l.br.AccountBalance(cash.ID)

		for i := 0; i < 5; i++ { // buy 5 kg of carrot
			_, _ = l.CreateTransaction(cash.ID, pubs[i].ID, time.Now(), 1.53, "1 bear")
		}

		closePartyBalance := l.br.AccountBalance(cash.ID)

		spendings := float64(closePartyBalance.Value-openPartyBalance.Value) / Million
		if spendings != -7.65 {
			t.Errorf("expected spendings for 5 bears: %+.2f, got: %+.2f", -7.65, spendings)
		}
	})

}

func TestListBalance(t *testing.T) {

	t.Parallel()

	br := CreateBalanceRegistry()

	aid := CreateID()
	bid := CreateID()
	tid := CreateID()

	val := func(v float64) int64 { return int64(v * Million) }

	// 3 redacts of the same balance:
	br.Add(
		Balance{ID: bid, Account: aid, Transaction: tid, Time: time.Now(), Value: val(1.11)})
	br.Add(
		Balance{ID: bid, Account: aid, Transaction: tid, Time: time.Now(), Value: val(1.12)})
	br.Add(
		Balance{ID: bid, Account: aid, Transaction: tid, Time: time.Now(), Value: val(3.15)})

	bid2 := CreateID()
	tid2 := CreateID()

	br.Add(
		Balance{ID: bid2, Account: aid, Transaction: tid2, Time: time.Now(), Value: val(1.76)})

	tBalances := make(map[ID]int)
	for _, b := range br.List() {
		tBalances[b.Transaction]++
	}

	if tBalances[tid] != 1 {
		t.Errorf("expected 1 last redaction of balance, found: %d", tBalances[tid])
	}

	if tBalances[tid2] != 1 {
		t.Errorf("expected 1 last redaction of balance[%s], found: %d", tid2, tBalances[tid2])
	}
}
