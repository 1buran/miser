package miser

import (
	"testing"
	"time"
)

func TestBalanceInit(t *testing.T) {

	t.Parallel()

	// Create repositories:
	ar := CreateAccountRegistry()
	tr := CreateTransactionRegistry()
	br := CreateBalanceRegistry()
	cr := CreateCurrencyRegistry()
	tg := CreateTagRegistry()
	tm := CreateTagsMapRegistry()

	// Create service:
	l := CreateLedger(ar, br, tr, cr, tg, tm)

	t.Run("zero", func(t *testing.T) {
		acc, err := l.CreateAccount("Deposit", Asset, "deposit account", "USD", time.Now(), 0.00)
		if err != nil {
			t.Fatal(err)
		}

		b := l.AccountBalance(acc.ID)
		if b == nil {
			t.Fatal("balance was not created during account creation")
		}

		if b.Value != 0 {
			t.Errorf("expected 0, %d found", b.Value)
		}
	})

	t.Run("float", func(t *testing.T) {
		acc, err := l.CreateAccount("Deposit", Asset, "deposit account", "USD", time.Now(), 123.78)
		if err != nil {
			t.Fatal(err)
		}

		if amount := l.AccountAmount(acc.ID); amount != 123.78 {
			t.Errorf("expected 123.78, got: %f", amount)
		}
	})
}

func TestBalanceChange(t *testing.T) {
	t.Parallel()

	// Create repositories:
	ar := CreateAccountRegistry()
	tr := CreateTransactionRegistry()
	br := CreateBalanceRegistry()
	cr := CreateCurrencyRegistry()
	tg := CreateTagRegistry()
	tm := CreateTagsMapRegistry()

	// Create service:
	l := CreateLedger(ar, br, tr, cr, tg, tm)

	t.Run("Expense", func(t *testing.T) {
		cash, err := l.CreateAccount("Cash", Asset, "wallet", "USD", time.Now(), 1555.12)
		if err != nil {
			t.Fatal(err)
		}

		market, err := l.CreateAccount("Market", Expense, "holiday market", "USD", time.Now(), 343.11)
		if err != nil {
			t.Fatal(err)
		}

		if openMarketBalance := l.AccountAmount(market.ID); openMarketBalance != 343.11 {
			t.Errorf("expected 344.64, got: %.2f", openMarketBalance)
		}

		_, err = l.CreateTransaction(cash.ID, market.ID, time.Now(), 1.53, "1kg carrot")
		if err != nil {
			t.Fatal(err)
		}

		if amount := l.AccountAmount(cash.ID); amount != 1553.59 {
			t.Errorf("expected 1553.59, got: %.2f", amount)
		}

		if updateMarketBalance := l.AccountAmount(market.ID); updateMarketBalance != 344.64 {
			t.Errorf("expected 344.64, got: %.2f", updateMarketBalance)
		}

	})

	t.Run("Earnings", func(t *testing.T) {
		var buyers []Account
		for i := 0; i < 5; i++ {
			acc, _ := l.CreateAccount("Cash", Asset, "wallet", "USD", time.Now(), 1555.12)
			buyers = append(buyers, *acc)
		}

		market, err := l.CreateAccount("Market", Expense, "holiday market", "USD", time.Now(), 343.11)
		if err != nil {
			t.Fatal(err)
		}

		openMarketBalance := l.AccountBalance(market.ID)
		t.Logf("open market balance: %#v", openMarketBalance)
		for i := 0; i < 5; i++ { // buy 5 kg of carrot
			_, _ = l.CreateTransaction(buyers[i].ID, market.ID, time.Now(), 1.53, "1kg carrot")
		}

		closeMarketBalance := l.AccountBalance(market.ID)
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
			acc, _ := l.CreateAccount("Shop", Asset, "LC Waikiki", "USD", time.Now(), 15.12)
			pubs = append(pubs, *acc)
		}

		cash, err := l.CreateAccount("Cash", Expense, "wallet", "USD", time.Now(), 343.11)
		if err != nil {
			t.Fatal(err)
		}

		openPartyBalance := l.AccountBalance(cash.ID)

		for i := 0; i < 5; i++ { // buy 5 kg of carrot
			_, _ = l.CreateTransaction(cash.ID, pubs[i].ID, time.Now(), 1.53, "1 bear")
		}

		closePartyBalance := l.AccountBalance(cash.ID)

		spendings := float64(closePartyBalance.Value-openPartyBalance.Value) / Million
		if spendings != -7.65 {
			t.Errorf("expected spendings for 5 bears: %+.2f, got: %+.2f", -7.65, spendings)
		}
	})

}

func TestBalanceList(t *testing.T) {

	t.Parallel()

	br := CreateBalanceRegistry()

	aid := CreateID()
	tid := CreateID()

	val := func(v float64) int64 { return int64(v * Million) }

	// 3 redacts of the same balance:
	br.Add(Balance{Account: aid, Transaction: tid, Value: val(1.11)})
	br.Add(Balance{Account: aid, Transaction: tid, Value: val(1.12)})
	br.Add(Balance{Account: aid, Transaction: tid, Value: val(3.15)})

	tid2 := CreateID()

	br.Add(Balance{Account: aid, Transaction: tid2, Value: val(1.76)})

	// count how many balances of transactions exist:
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

func TestBalanceUpdate(t *testing.T) {

	t.Parallel()

	// Create repositories:
	ar := CreateAccountRegistry()
	tr := CreateTransactionRegistry()
	br := CreateBalanceRegistry()
	cr := CreateCurrencyRegistry()
	tg := CreateTagRegistry()
	tm := CreateTagsMapRegistry()

	// Create service:
	l := CreateLedger(ar, br, tr, cr, tg, tm)

	t.Run("Linear", func(t *testing.T) {
		openedAt := time.Date(2024, time.January, 1, 8, 30, 0, 0, time.UTC)
		wallet, err := l.CreateAccount("Cash", Asset, "wallet", "USD", openedAt, 1555.12)
		if err != nil {
			t.Fatal(err)
		}

		bazaar, err := l.CreateAccount("Bazaar", Expense, "sunday bazaar", "USD", openedAt, 0.50)
		if err != nil {
			t.Fatal(err)
		}

		shopingTime := openedAt.Add(2 * time.Hour)
		if _, err := l.CreateTransaction(wallet.ID, bazaar.ID, shopingTime, 5.35, "oranges"); err != nil {
			t.Fatal(err)
		}

		if l.AccountAmount(wallet.ID) != 1549.77 {
			t.Errorf("expected 1549.77 in wallet, found: %.2f", l.AccountAmount(wallet.ID))
		}

		if l.AccountAmount(bazaar.ID) != 5.85 {
			t.Errorf("expected 5.85 in wallet, found: %.2f", l.AccountAmount(bazaar.ID))
		}

		shopingTime2 := shopingTime.Add(30 * time.Minute)
		if _, err := l.CreateTransaction(wallet.ID, bazaar.ID, shopingTime2, 2.13, "carrot"); err != nil {
			t.Fatal(err)
		}

		if l.AccountAmount(wallet.ID) != 1547.64 {
			t.Errorf("expected 1547.64 in wallet, found: %.2f", l.AccountAmount(wallet.ID))
		}

		if l.AccountAmount(bazaar.ID) != 7.98 {
			t.Errorf("expected 7.98 in wallet, found: %.2f", l.AccountAmount(bazaar.ID))
		}
	})

	t.Run("Intermediate", func(t *testing.T) {
		openedAt := time.Date(2024, time.October, 1, 15, 30, 0, 0, time.UTC)
		wallet, err := l.CreateAccount("Cash", Asset, "wallet", "USD", openedAt, 200.37)
		if err != nil {
			t.Fatal(err)
		}

		bazaar, err := l.CreateAccount("Bazaar", Expense, "sunday bazaar", "USD", openedAt, 0.50)
		if err != nil {
			t.Fatal(err)
		}

		// fist transaction at 20 Oct 15:30:
		dt1 := time.Date(2024, time.October, 20, 15, 30, 0, 0, time.UTC)
		tr1, err := l.CreateTransaction(wallet.ID, bazaar.ID, dt1, 2.13, "carrot")
		if err != nil {
			t.Fatal(err)
		}

		// check account balance after first committed transaction:
		bw1 := l.br.TransactionBalance(wallet.ID, tr1.ID)
		if bw1 == nil {
			t.Fatalf("balance not found, account: %s, transaction: %s", wallet.ID, tr1.ID)
		}
		if bw1.Amount() != 200.37-2.13 {
			t.Errorf("expected balance of wallet: %.2f, found: %.2f", 200.37-2.13, bw1.Amount())
		}
		bz1 := l.br.TransactionBalance(bazaar.ID, tr1.ID)
		if bz1 == nil {
			t.Fatalf("exoected bazaar income: %.2f, found: %.2f", 0.50+2.13, bz1.Amount())
		}

		// second transaction at 20 Oct 22:30:
		dt2 := time.Date(2024, time.October, 20, 22, 30, 0, 0, time.UTC)
		tr2, err := l.CreateTransaction(wallet.ID, bazaar.ID, dt2, 5.17, "oranges")
		if err != nil {
			t.Fatal(err)
		}

		bw2 := l.br.TransactionBalance(wallet.ID, tr2.ID)
		if bw2 == nil {
			t.Fatalf("balance not found, account: %s, transaction: %s", wallet.ID, tr2.ID)
		}
		if bw2.Amount() != 200.37-2.13-5.17 {
			t.Errorf("expected balance of wallet: %.2f, found: %.2f", 200.37-2.13-5.17, bw2.Amount())
		}
		bz2 := l.br.TransactionBalance(bazaar.ID, tr2.ID)
		if bz2 == nil {
			t.Fatalf("exoected bazaar income: %.2f, found: %.2f", 0.50+2.13+5.17, bz2.Amount())
		}

		// third transaction, add forgotten expense, cheese after end of workday,
		// between first and second transactions:
		dt3 := time.Date(2024, time.October, 20, 17, 30, 0, 0, time.UTC)
		tr3, err := l.CreateTransaction(wallet.ID, bazaar.ID, dt3, 150, "cheese")
		if err != nil {
			t.Fatal(err)
		}

		// check account balance after third committed transaction,
		// balance at time of first transaction should not changed,
		// balance at time of second transaction affected:
		bw1 = l.br.TransactionBalance(wallet.ID, tr1.ID)
		if bw1 == nil {
			t.Fatalf("balance not found, account: %s, transaction: %s", wallet.ID, tr1.ID)
		}
		if bw1.Amount() != 200.37-2.13 {
			t.Errorf("expected balance of wallet: %.2f, found: %.2f", 200.37-2.13, bw1.Amount())
		}

		bw3 := l.br.TransactionBalance(wallet.ID, tr3.ID)
		if bw3 == nil {
			t.Fatalf("balance not found, account: %s, transaction: %s", wallet.ID, tr3.ID)
		}
		if bw3.Amount() != 200.37-2.13-150 {
			t.Errorf("expected balance of wallet: %.2f, found: %.2f", 200.37-2.13-150, bw3.Amount())
		}

		bw2 = l.br.TransactionBalance(wallet.ID, tr2.ID)
		if bw2 == nil {
			t.Fatalf("balance not found, account: %s, transaction: %s", wallet.ID, tr2.ID)
		}
		if bw2.Amount() != 200.37-2.13-5.17-150 {
			t.Errorf("expected balance of wallet: %.2f, found: %.2f", 200.37-2.13-5.17-150, bw2.Amount())
		}

		if l.AccountAmount(wallet.ID) != 200.37-2.13-5.17-150 {
			t.Errorf("expected balance of account: %.2f, got: %.2f", 200.37-2.13-5.17-150, l.AccountAmount(wallet.ID))
		}

		bz1 = l.br.TransactionBalance(bazaar.ID, tr1.ID)
		if bz1 == nil {
			t.Fatalf("exoected bazaar income: %.2f, found: %.2f", 0.50+2.13, bz1.Amount())
		}

		bz3 := l.br.TransactionBalance(bazaar.ID, tr3.ID)
		if bz1 == nil {
			t.Fatalf("exoected bazaar income: %.2f, found: %.2f", 0.50+2.13+150, bz3.Amount())
		}

		// balance of second transaction of bazaar account was affected:
		bz2 = l.br.TransactionBalance(bazaar.ID, tr2.ID)
		if bz2 == nil {
			t.Fatalf("exoected bazaar income: %.2f, found: %.2f", 0.50+2.13+5.17+150, bz2.Amount())
		}

		if l.AccountAmount(bazaar.ID) != 0.50+2.13+5.17+150 {
			t.Errorf("expected balance of account: %.2f, got: %.2f", 0.50+2.13+5.17+150, l.AccountAmount(bazaar.ID))
		}

	})
}
