package miser

import (
	"testing"
	"time"
)

func TestInitBalance(t *testing.T) {

	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		acc, err := CreateAccount("Deposit", Asset, "deposit account", "USD", 0.00)
		if err != nil {
			t.Fatal(err)
		}

		b := Balances.AccountBalance(acc.ID)
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
		acc, err := CreateAccount("Deposit", Asset, "deposit account", "USD", 123.78)
		if err != nil {
			t.Fatal(err)
		}

		if amount := Balances.AccountValue(acc.ID); amount != 123.78 {
			t.Errorf("expected 123.78, got: %f", amount)
		}
	})
}

func TestChangeBalance(t *testing.T) {
	t.Parallel()

	t.Run("Expense", func(t *testing.T) {
		cash, err := CreateAccount("Cash", Asset, "wallet", "USD", 1555.12)
		if err != nil {
			t.Fatal(err)
		}

		market, err := CreateAccount("Market", Expense, "holiday market", "USD", 343.11)
		if err != nil {
			t.Fatal(err)
		}

		if openMarketBalance := Balances.AccountValue(market.ID); openMarketBalance != 343.11 {
			t.Errorf("expected 344.64, got: %.2f", openMarketBalance)
		}

		_, err = CreateTransation(cash.ID, market.ID, time.Now(), 1.53, "1kg carrot")
		if err != nil {
			t.Fatal(err)
		}

		if amount := Balances.AccountValue(cash.ID); amount != 1553.59 {
			t.Errorf("expected 1553.59, got: %.2f", amount)
		}

		if updateMarketBalance := Balances.AccountValue(market.ID); updateMarketBalance != 344.64 {
			t.Errorf("expected 344.64, got: %.2f", updateMarketBalance)
		}

	})

	t.Run("Earnings", func(t *testing.T) {
		var buyers []Account
		for i := 0; i < 5; i++ {
			acc, _ := CreateAccount("Cash", Asset, "wallet", "USD", 1555.12)
			buyers = append(buyers, *acc)
		}

		market, err := CreateAccount("Market", Expense, "holiday market", "USD", 343.11)
		if err != nil {
			t.Fatal(err)
		}

		openMarketBalance := Balances.AccountBalance(market.ID)
		t.Logf("open market balance: %#v", openMarketBalance)
		for i := 0; i < 5; i++ { // buy 5 kg of carrot
			_, _ = CreateTransation(buyers[i].ID, market.ID, time.Now(), 1.53, "1kg carrot")
		}

		closeMarketBalance := Balances.AccountBalance(market.ID)
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
			acc, _ := CreateAccount("Shop", Asset, "LC Waikiki", "USD", 15.12)
			pubs = append(pubs, *acc)
		}

		cash, err := CreateAccount("Cash", Expense, "wallet", "USD", 343.11)
		if err != nil {
			t.Fatal(err)
		}

		openPartyBalance := Balances.AccountBalance(cash.ID)

		for i := 0; i < 5; i++ { // buy 5 kg of carrot
			_, _ = CreateTransation(cash.ID, pubs[i].ID, time.Now(), 1.53, "1 bear")
		}

		closePartyBalance := Balances.AccountBalance(cash.ID)

		spendings := float64(closePartyBalance.Value-openPartyBalance.Value) / Million
		if spendings != -7.65 {
			t.Errorf("expected spendings for 5 bears: %+.2f, got: %+.2f", -7.65, spendings)
		}
	})

}

func TestListBalance(t *testing.T) {

	t.Parallel()

	aid := CreateID()
	bid := CreateID()
	tid := CreateID()

	val := func(v float64) int64 { return int64(v * Million) }

	lenBalances := len(Balances.Items)

	// 3 redacts of the same balance:
	Balances.Add(
		Balance{ID: bid, Account: aid, Transaction: tid, Time: time.Now(), Value: val(1.11)})
	Balances.Add(
		Balance{ID: bid, Account: aid, Transaction: tid, Time: time.Now(), Value: val(1.12)})
	Balances.Add(
		Balance{ID: bid, Account: aid, Transaction: tid, Time: time.Now(), Value: val(3.15)})

	bid2 := CreateID()
	Balances.Add(
		Balance{ID: bid2, Account: aid, Transaction: tid, Time: time.Now(), Value: val(31.76)})

	balances := Balances.List()
	if len(balances) != 2+lenBalances {
		t.Errorf("expected balances length: 2, got: %d", len(balances))
	}
	t.Logf("%#v", balances)
}
