package entity

import (
	"fmt"
	"time"
)

const TransactionTypeTransfer = "Transfer"
const TransactionTypeDeposit = "Deposit"
const TransactionTypeWithdraw = "Withdraw"

type Transaction struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	SourceID uint
	// SourceAccount      Account
	SourceAmount float64

	DestinationID uint
	// DestinationAccount Account
	DestinationAmount float64

	Type string

	Rate     float64
	RateType string
}

func (Transaction) TableName() string {
	return "trading_bot.transactions"
}

func (t Transaction) String() string {
	sourceID := fmt.Sprintf("%d", uint64(t.SourceID))
	destinationID := fmt.Sprintf("%d", uint64(t.DestinationID))
	sourceAmount := fmt.Sprintf("%f", t.SourceAmount)
	destinationAmount := fmt.Sprintf("%f", t.DestinationAmount)

	return t.Type + " " + sourceID + " " + sourceAmount + " " + destinationID + " " + destinationAmount
}
