package main

const BID = "bid"
const OFFER = "offer"

var configs = map[string]quoteConfig{}

type entity struct {
	size      float64
	direction string
	price     float64
}

type quoteConfig struct {
	interval int64
	entities []entity
}

func LoadDefaultQuoteConfig() {
	configs["EUR/USD_TOM"] = quoteConfig{
		interval: 10000,
		entities: []entity{
			{size: 1000, direction: BID, price: 1.1},
			{size: 1000, direction: OFFER, price: 1.2},
			{size: 1000000, direction: BID, price: 1.05},
		},
	}
}

func GetQuoteConfig(symbol string) *quoteConfig {
	res, ok := configs[symbol]
	if ok {
		return &res
	}
	return nil
}
