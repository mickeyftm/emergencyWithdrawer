package config

import (
	"encoding/json"
	"errors"
	"log"
)

// https://github.com/shibukawa/configdir

func init() {
	err := json.Unmarshal([]byte(defaultConfig), &conf)
	if err != nil {
		log.Panicln(err)
	}
}

var (
	conf       map[string]*config
	activeConf *activeConfig
)

var (
	ErrNoConfFound  = errors.New("no config for the selected network provided")
	ErrNoActiveConf = errors.New("no config loaded")
)

type activeConfig struct {
	Name string
	config
}

type config struct {
	Endpoint   string  `json:"endpoint"`
	GasPrice   float64 `json:"gasPrice"`
	GasLimit   float64 `json:"gasLimit"`
	PrivateKey string  `json:"-"`
}

func SetActiveConf(net string) error {
	cfg, ok := conf[net]
	if !ok {
		return ErrNoConfFound
	}

	activeConf = &activeConfig{
		Name:   net,
		config: *cfg,
	}
	return nil
}

func GetActiveConf() (*activeConfig, error) {
	cfg := activeConf
	if cfg != nil {
		return cfg, nil
	}
	return nil, ErrNoActiveConf
}

const defaultConfig = `
{
	"eth": {
		"endpoint": "https://main-light.eth.linkpool.io",
		"gasPrice": 0,
		"gasLimit": 100000
	},
	"bsc": {
		"endpoint": "https://bsc-dataseed.binance.org",
		"gasPrice": 0,
		"gasLimit": 100000
	},
	"matic": {
		"endpoint": "https://rpc-mainnet.matic.quiknode.pro",
		"gasPrice": 0,
		"gasLimit": 100000
	},
	"ftm": {
		"endpoint": "https://rpcapi.fantom.network",
		"gasPrice": 0,
		"gasLimit": 100000
	},
	"kcc": {
		"endpoint": "https://rpc-mainnet.kcc.network",
		"gasPrice": 0,
		"gasLimit": 100000
	}
}
`
