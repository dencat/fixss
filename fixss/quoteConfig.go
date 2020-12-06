package fixss

import (
	"github.com/quickfixgo/enum"
	"sort"
)

const BID = "bid"
const OFFER = "offer"

var quoteConfigs = map[string]QuoteConfig{}

type entity struct {
	Size      float64 `json:"size"`
	Direction string  `json:"direction"`
	MinPrice  float64 `json:"minPrice"`
	MaxPrice  float64 `json:"maxPrice"`
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
			{Size: 1000, Direction: BID, MinPrice: 1.1, MaxPrice: 1.11},
			{Size: 1000, Direction: OFFER, MinPrice: 1.2, MaxPrice: 1.21},
			{Size: 1000000, Direction: BID, MinPrice: 1.05, MaxPrice: 1.07},
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
	sort.Slice(quoteConfig.Entities, func(i, j int) bool {
		if quoteConfig.Entities[i].Size == quoteConfig.Entities[j].Size {
			return quoteConfig.Entities[i].Direction < quoteConfig.Entities[j].Direction
		}

		return quoteConfig.Entities[i].Size < quoteConfig.Entities[j].Size
	})
}

func GetMarketPrice(symbol string, side enum.Side, size float64) *float64 {
	var res *float64 = nil
	if config, ok := quoteConfigs[symbol]; ok {
		for _, entity := range config.Entities {
			if size > entity.Size {
				continue
			}
			if entity.Direction == BID && side == enum.Side_BUY {
				continue
			}
			if entity.Direction == OFFER && side == enum.Side_SELL {
				continue
			}
			return &entity.MinPrice
		}
	}
	return res
}
