package fixss

import (
	"github.com/shopspring/decimal"
	"math/rand"
)

type Quote struct {
	price     decimal.Decimal
	Size      float64
	Direction string
}

func GetNextQuotes(symbol string) ([]Quote, int64) {
	quoteConfig := GetQuoteConfig(symbol)
	if quoteConfig == nil {
		return nil, 0
	}

	quotes := make([]Quote, 0)
	for _, entity := range quoteConfig.Entities {
		quote := Quote{
			price:     decimal.NewFromFloat(entity.MinPrice + rand.Float64()*(entity.MaxPrice-entity.MinPrice)),
			Size:      entity.Size,
			Direction: entity.Direction,
		}
		quotes = append(quotes, quote)
	}

	return quotes, quoteConfig.Interval
}
