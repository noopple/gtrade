package entity

import "gorm.io/gorm"

type Deal struct {
	gorm.Model

	PairID uint
	Pair Pair

	Orders []Order
}

func (Deal) TableName() string {
	return "trading_bot.deals"
}
