package main

const BID = "bid"
const OFFER = "offer"

var configs = map[string]QuoteConfig{}

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
	configs["EUR/USD_TOM"] = QuoteConfig{
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
	res, ok := configs[symbol]
	if ok {
		return &res
	}
	return nil
}

func SetQuoteConfig(quoteConfig QuoteConfig) {
	configs[quoteConfig.Symbol] = quoteConfig
}
