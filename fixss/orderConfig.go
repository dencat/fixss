package fixss

type StrategyType string

var orderConfigs = map[string]OrderConfig{}

const DEFAULT = "default"

const (
	Accept StrategyType = "accept"
	Reject StrategyType = "reject"
)

type OrderConfig struct {
	Symbol   string       `json:"symbol"`
	Strategy StrategyType `json:"strategy"`
}

func GetOrderConfig(symbol string) OrderConfig {
	if config, ok := orderConfigs[symbol]; ok {
		return config
	}
	if config, ok := orderConfigs[DEFAULT]; ok {
		return config
	}

	return OrderConfig{
		Strategy: Accept,
	}
}

func SetOrderConfig(orderConfig OrderConfig) {
	Log.Infof("Set order config %s", orderConfig)
	orderConfigs[orderConfig.Symbol] = orderConfig
}
