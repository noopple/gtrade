package entity

import "gorm.io/gorm"

const InstrumentTypeCurrency = "Currency"
const InstrumentTypeStock = "Stock"

type Instrument struct {
	gorm.Model

	Symbol string
	Type string
}

func (Instrument) TableName() string {
	return "trading_bot.instruments"
}
