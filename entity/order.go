package entity

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	DealID uint
	Deal Deal
	Price float64
	Quantity float64
	Operation string
	TransactionID *uint
	Transaction *Transaction

}

func (Order) TableName() string {
	return "trading_bot.orders"
}
