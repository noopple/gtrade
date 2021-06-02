package terminal

import (
	"errors"

	"gorm.io/gorm"
	"gtrade/database"
	"gtrade/entity"
)

type Terminal struct {
	Symbol string

	db *gorm.DB
}

func Get(symbol string) *Terminal {
	db, _ := database.GetConnection()
	return &Terminal{db: db, Symbol: symbol}
}

func (t Terminal) NewOrder(symbol string, price, quantity float64, operation string) entity.Order {
	var pair entity.Pair
	err := t.db.
		Joins("BaseInstrument").
		Joins("QuoteInstrument").
		Where("\"BaseInstrument\".symbol = ? AND \"QuoteInstrument\".symbol = ?", t.Symbol, symbol).
		Take(&pair).
		Error

	if err != nil {
		return entity.Order{}
	}

	deal, err := t.getOpenDeal(pair.ID)

	if err != nil {
		deal = t.newDeal(pair.ID)
	}

	order := entity.Order{
		DealID: deal.ID,
		Price: price,
		Quantity: quantity,
		Operation: operation,
	}

	t.db.Create(&order)

	return order
}

func (t Terminal) ProcessOrder(order entity.Order) {
	var deal entity.Deal
	var sourceAccount entity.Account
	var destinationAccount entity.Account

	t.db.Where(order.DealID).Joins("Pair").First(&deal)
	t.db.Where("instrument_id = ?", deal.Pair.BaseInstrumentID).Take(&sourceAccount)
	t.db.Where("instrument_id = ?", deal.Pair.QuoteInstrumentID).Take(&destinationAccount)

	transaction := entity.Transaction{
		Type: entity.TransactionTypeTransfer,
		Rate: order.Price,
	}

	var rateType string
	if order.Operation == "sell" {
		rateType = "divisor"
		transaction.SourceID = destinationAccount.ID
		transaction.DestinationID = sourceAccount.ID
		transaction.SourceAmount = order.Quantity
		transaction.DestinationAmount = order.Quantity * order.Price
	} else {
		rateType = "multiplier"
		transaction.SourceID = sourceAccount.ID
		transaction.DestinationID = destinationAccount.ID
		transaction.SourceAmount = order.Quantity * order.Price
		transaction.DestinationAmount = order.Quantity
	}

	transaction.RateType = rateType

	t.db.Create(&transaction)

	order.TransactionID = &transaction.ID

	t.db.Save(&order)
}

func (t Terminal) newDeal(pairID uint) entity.Deal {
	deal := entity.Deal{
		PairID: pairID,
	}
	t.db.Create(&deal)

	return deal
}

func (t Terminal) getOpenDeal(pairID uint) (entity.Deal, error) {
	var dealID uint
	var deal entity.Deal

	t.db.Model(&entity.Deal{}).
		Joins("Orders").
		Where("pair_id = ?", pairID).
		Group("deals.id").
		Having("SUM(CASE WHEN operation = 'buy' THEN quantity ELSE quantity END) <> 0 OR SUM(quantity) IS NULL").
		Pluck("deals.id", &dealID)

	if dealID == 0 {
		return entity.Deal{}, errors.New("open deal is not found")
	}

	t.db.Take(&deal, dealID)

	return deal, nil
}

