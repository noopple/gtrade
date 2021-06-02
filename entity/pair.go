package entity

type Pair struct {
	// gorm.Model

	// PairID uint
	ID uint
	BaseInstrumentID uint
	BaseInstrument Instrument `gorm:"foreignKey:BaseInstrumentID"`

	QuoteInstrumentID uint
	QuoteInstrument Instrument `gorm:"foreignKey:QuoteInstrumentID"`
}

func (Pair) TableName() string {
	return "trading_bot.pairs"
}
