package config

import (
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"sync"
)

// https://github.com/shibukawa/configdir

func init() {
	err := json.Unmarshal([]byte(defaultConfig), &conf)
	if err != nil {
		log.Panicln(err)
	}
}

var (
	mu         sync.RWMutex
	conf       map[string]*config
	activeConf *ActiveConfig
)

var (
	ErrNoConfFound = errors.New("no config for the selected network provided")
	// return this error when no active config is set
	ErrNoActiveConf = errors.New("no config loaded")
)

type ActiveConfig struct {
	Name string
	config
}

type config struct {
	Endpoint   string   `json:"endpoint"`
	GasPrice   *big.Int `json:"gasPrice"`
	GasLimit   uint64   `json:"gasLimit"`
	TxExplorer string   `json:"txExplorer"`
}

func SetActiveConf(net string) error {
	mu.Lock()
	defer mu.Unlock()

	cfg, ok := conf[net]
	if !ok {
		return ErrNoConfFound
	}

	activeConf = &ActiveConfig{
		Name:   net,
		config: *cfg,
	}
	return nil
}

func GetActiveConf() (*ActiveConfig, error) {
	mu.RLock()
	defer mu.RUnlock()

	cfg := activeConf
	if cfg != nil {
		return cfg, nil
	}
	return nil, ErrNoActiveConf
}

func SetActiveEndpoint(endpoint string) error {
	mu.Lock()
	defer mu.Unlock()

	if activeConf == nil {
		activeConf = &ActiveConfig{}
	}
	activeConf.config.Endpoint = endpoint
	return nil
}

func SetActiveGasPrice(gasPrice *big.Int) error {
	mu.Lock()
	defer mu.Unlock()

	if activeConf == nil {
		activeConf = &ActiveConfig{}
	}
	activeConf.config.GasPrice = gasPrice
	return nil
}

func SetActiveGasLimit(gasLimit uint64) error {
	mu.Lock()
	defer mu.Unlock()

	if activeConf == nil {
		activeConf = &ActiveConfig{}
	}
	activeConf.config.GasLimit = gasLimit
	return nil
}

const defaultConfig = `
{
	"eth": {
		"endpoint": "https://main-light.eth.linkpool.io",
		"gasPrice": 0,
		"gasLimit": 10000000,
		"txExplorer": "https://etherscan.io/tx/"
	},
	"bsc": {
		"endpoint": "https://bsc-dataseed.binance.org",
		"gasPrice": 0,
		"gasLimit": 10000000,
		"txExplorer": "https://bscscan.com/tx/"
	},
	"matic": {
		"endpoint": "https://rpc-mainnet.matic.quiknode.pro",
		"gasPrice": 0,
		"gasLimit": 10000000,
		"txExplorer": "https://polygonscan.com/tx/"
	},
	"ftm": {
		"endpoint": "https://rpcapi.fantom.network",
		"gasPrice": 0,
		"gasLimit": 10000000,
		"txExplorer": "https://ftmscan.com/tx/"
	},
	"kcc": {
		"endpoint": "https://rpc-mainnet.kcc.network",
		"gasPrice": 0,
		"gasLimit": 10000000,
		"txExplorer": "https://explorer.kcc.io/en/tx/"
	}
}
`
