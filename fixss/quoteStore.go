package fixss

import (
	"github.com/quickfixgo/enum"
	"github.com/shopspring/decimal"
	"math/rand"
)

type Quote struct {
	price     decimal.Decimal
	Size      float64
	Direction string
}

var quoteStore = map[string][]Quote{}

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
	quoteStore[symbol] = quotes

	return quotes, quoteConfig.Interval
}

func GetMarketQuote(symbol string, side enum.Side, size float64) *Quote {
	var res *Quote = nil
	if quotes, ok := quoteStore[symbol]; ok {
		for _, quote := range quotes {
			if size > quote.Size {
				continue
			}
			if quote.Direction == BID && side == enum.Side_BUY {
				continue
			}
			if quote.Direction == OFFER && side == enum.Side_SELL {
				continue
			}
			return &quote
		}
	}
	return res
}
