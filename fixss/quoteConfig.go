package fixss

import "github.com/quickfixgo/enum"

const BID = "bid"
const OFFER = "offer"

var quoteConfigs = map[string]QuoteConfig{}

type entity struct {
	Size      float64 `json:"size"`
	Direction string  `json:"direction"`
	Price     float64 `json:"price"`
}

type QuoteConfig struct {
	Symbol   string   `json:"symbol"`
	Interval int64    `json:"interval"`
	Entities []entity `json:"entities"`
}

func LoadDefaultQuoteConfig() {
	quoteConfigs["EUR/USD_TOM"] = QuoteConfig{
		Symbol:   "EUR/USD_TOM",
		Interval: 10000,
		Entities: []entity{
			{Size: 1000, Direction: BID, Price: 1.1},
			{Size: 1000, Direction: OFFER, Price: 1.2},
			{Size: 1000000, Direction: BID, Price: 1.05},
		},
	}
}

func GetQuoteConfig(symbol string) *QuoteConfig {
	res, ok := quoteConfigs[symbol]
	if ok {
		return &res
	}
	return nil
}

func SetQuoteConfig(quoteConfig QuoteConfig) {
	quoteConfigs[quoteConfig.Symbol] = quoteConfig
}

func GetMarketPrice(symbol string, side enum.Side, size float64) *float64 {
	var res *float64 = nil
	if config, ok := quoteConfigs[symbol]; ok {
		for _, entity := range config.Entities {
			if entity.Direction == BID && side == enum.Side_BUY {
				continue
			}
			if entity.Direction == OFFER && side == enum.Side_SELL {
				continue
			}
		}
	}
	return res
}
