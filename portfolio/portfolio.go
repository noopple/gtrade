package portfolio

import (
	"fmt"

	"github.com/piquette/finance-go/quote"
	"gorm.io/gorm"
	"gtrade/database"
	"gtrade/entity"
)

type Portfolio struct {
	IsEmpty  bool
	Account  entity.Account
	Balances []*Balance

	// Deposits []*Balance
	// Orders []entity.Order
	// Currencies []Balance
	// Stocks []Balance

	Symbol string

	Income float64
	Total  float64

	rates   map[string]float64
	transactions []entity.Transaction

	db      *gorm.DB
}

type Balance struct {
	Account entity.Account
	Value   float64
}

func Get(symbol string) Portfolio {
	var portfolio Portfolio

	portfolio.rates = make(map[string]float64)
	portfolio.db, _ = database.GetConnection()

	portfolio.Symbol = symbol

	portfolio.getTransactions()
	portfolio.getAccounts()

	portfolio.getBalances()

	portfolio.getRates()

	portfolio.calculateTotal()
	portfolio.calculateIncome()

	return portfolio
}

func (p *Portfolio) calculateTotal() {
	for _, balance := range p.Balances {
		if balance.Account.Instrument.Symbol == p.Symbol {
			p.Total += balance.Value
		} else {
			p.Total += balance.Value * p.rates[balance.Account.Instrument.Symbol]
		}
	}
}

func (p *Portfolio) getAccounts() {
	var accounts []entity.Account

	p.db.Joins("Instrument").Find(&accounts)

	if len(accounts) == 0 {
		p.IsEmpty = true
		return
	}

	for _, account := range accounts {
		p.Balances = append(p.Balances, &Balance{Account: account, Value: 0})

		if account.Instrument.Symbol == p.Symbol {
			p.Account = account
		}
	}
}

func (p *Portfolio) getTransactions() {
	p.db.Find(&p.transactions)
}

func (p *Portfolio) getStocks() []string {
	var stocks []string
	for _, balance := range p.Balances {
		if balance.Account.Instrument.Type == entity.InstrumentTypeStock {
			stocks = append(stocks, balance.Account.Instrument.Symbol)
		}
	}
	return stocks
}

func (p *Portfolio) getBalances() {
	balances := make(map[uint]*Balance)

	for _, balance := range p.Balances {
		balances[balance.Account.ID] = balance
	}

	for _, transaction := range p.transactions {
		if transaction.DestinationID > 0 {
			balances[transaction.DestinationID].Value += transaction.DestinationAmount
		}
		if transaction.SourceID > 0 {
			balances[transaction.SourceID].Value -= transaction.SourceAmount
		}
	}

	for index, balance := range p.Balances {
		p.Balances[index].Value = balances[balance.Account.ID].Value
	}
}

func (p Portfolio) getRates() {
	symbols := p.getStocks()
	quotes := quote.List(symbols)

	for quotes.Next() {
		p.rates[quotes.Quote().Symbol] = quotes.Quote().RegularMarketPrice
	}
}

func (p Portfolio) String() string {
	var description string
	for _, balance := range p.Balances {
		if balance.Value == 0 {
			continue
		}

		symbol := balance.Account.Instrument.Symbol
		amount := fmt.Sprintf("%f", balance.Value)
		rate   := fmt.Sprintf("%f", p.rates[symbol])

		row := symbol + "(" + rate + "): " + amount + "\n"
		if balance.Account.Instrument.Type == entity.InstrumentTypeCurrency {
			row = symbol + ": " + amount + "\n"
			description =  row + description
		} else {
			description = description + row
		}
	}

	return description + "Total: " + fmt.Sprintf("%f", p.Total) + " (" + fmt.Sprintf("%f", p.Income) + ")"
}

func (p *Portfolio) calculateIncome() {
	for _, transaction := range p.transactions {
		if transaction.SourceID == p.Account.ID && transaction.Type == entity.TransactionTypeWithdraw {
			p.Income -= transaction.SourceAmount
		}

		if transaction.DestinationID == p.Account.ID && transaction.Type == entity.TransactionTypeDeposit {
			p.Income += transaction.DestinationAmount
		}
	}
	p.Income = p.Total - p.Income
}