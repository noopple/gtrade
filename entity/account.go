package entity

import "gorm.io/gorm"

type Account struct {
	gorm.Model

	InstrumentID uint

	Instrument Instrument
}

func (Account) TableName() string {
	return "trading_bot.accounts"
}