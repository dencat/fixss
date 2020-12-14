package fixss

import (
	"encoding/json"
	"errors"
	log "github.com/jeanphorn/log4go"
	"github.com/quickfixgo/enum"
	"io/ioutil"
	"os"
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

func LoadDefaultQuoteConfig(config *Config) error {
	if len(config.Quote.Config) > 0 {
		if !FileExists(config.Quote.Config) {
			return errors.New("File " + config.Quote.Config + " not found")
		}
		jsonFile, err := os.Open(config.Quote.Config)
		if err != nil {
			return err
		}
		defer jsonFile.Close()
		data, _ := ioutil.ReadAll(jsonFile)
		var quoteConfigs []QuoteConfig
		err = json.Unmarshal(data, &quoteConfigs)
		if err != nil {
			return err
		}

		log.LOGGER("app").Info("Load default quote config from: %s", config.Quote.Config)

		for _, c := range quoteConfigs {
			SetQuoteConfig(c)
		}
	}
	return nil
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
