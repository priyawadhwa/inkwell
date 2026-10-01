package main

import (
	"fmt"
	"strings"
	"time"
)

// Cents avoids floating point for money.
type Cents int64

func (c Cents) String() string {
	sign := ""
	if c < 0 {
		sign = "-"
		c = -c
	}
	whole := fmt.Sprintf("%d", c/100)
	var b strings.Builder
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	return fmt.Sprintf("%s$%s.%02d", sign, b.String(), c%100)
}

type Account struct {
	Name    string
	Number  string
	Balance Cents
	Note    string
}

type Transaction struct {
	When     time.Time
	Merchant string
	Category string
	Amount   Cents
}

type Card struct {
	Name   string
	Last4  string
	Frozen bool
}

type Customer struct {
	FirstName    string
	Accounts     []Account
	Transactions []Transaction
	Cards        []Card
}

func (c Customer) NetWorth() Cents {
	var total Cents
	for _, a := range c.Accounts {
		total += a.Balance
	}
	return total
}

// demoCustomer is the one customer the demo serves. Linky is Chainguard's
// octopus mascot.
func demoCustomer(now time.Time) Customer {
	day := func(n int) time.Time { return now.AddDate(0, 0, -n) }
	return Customer{
		FirstName: "Linky",
		Accounts: []Account{
			{Name: "Everyday Checking", Number: "••8808", Balance: 1248022, Note: "Eight arms, one account"},
			{Name: "Reef Savings", Number: "••4120", Balance: 3861500, Note: "Round-ups go to the reef"},
			{Name: "Tide Pool Brokerage", Number: "••0007", Balance: 9120419, Note: "Diversified across all seven seas"},
		},
		Transactions: []Transaction{
			{When: day(0), Merchant: "Coral Café", Category: "Coffee", Amount: -675},
			{When: day(0), Merchant: "Kelp & Co. Grocers", Category: "Groceries", Amount: -8412},
			{When: day(1), Merchant: "Tidepool Transit", Category: "Transport", Amount: -275},
			{When: day(2), Merchant: "Payroll — Chainguard", Category: "Income", Amount: 512000},
			{When: day(3), Merchant: "Squid Ink Pasta Bar", Category: "Dining", Amount: -6230},
			{When: day(4), Merchant: "Deep Sea Cable Co.", Category: "Internet", Amount: -7999},
			{When: day(6), Merchant: "Sunken Treasure Antiques", Category: "Shopping", Amount: -24500},
		},
		Cards: []Card{
			{Name: "Inkwell Debit", Last4: "8808"},
			{Name: "Abyss Rewards Visa", Last4: "2045"},
		},
	}
}
