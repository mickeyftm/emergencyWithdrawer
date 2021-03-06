package client

import (
	"errors"
	"sync"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jon4hz/emergenyWithdrawer/internal/config"
	"github.com/jon4hz/web3-multicall-go/multicall"
)

var (
	cm              sync.Mutex // protects the clients
	Client          *ethclient.Client
	ETHClient       *ethrpc.ETH
	MulticallClient multicall.Multicall

	multicallContract = map[string]string{
		"eth":   "0x5eb3fa2dfecdde21c950813c665e9364fa609bd2",
		"bsc":   "0x6Cf63cC81660Dd174A49e0C61A1f916456Ee1471",
		"matic": "0x8a233a018a2e123c0D96435CF99c8e65648b429F",
		"ftm":   "0x08AB4aa09F43cF2D45046870170dd75AE6FBa306",
		"kcc":   "0x08AB4aa09F43cF2D45046870170dd75AE6FBa306",
	}
)

var (
	ErrNoMulticall = errors.New("network has no multicall contract")
)

func Load(net string) error {

	err := config.SetActiveConf(net)
	if err != nil {
		return err
	}
	cfg, err := config.GetActiveConf()
	if err != nil {
		return err
	}

	cm.Lock()
	defer cm.Unlock()

	Client, err = ethclient.Dial(cfg.Endpoint)
	if err != nil {
		return err
	}

	ETHClient, err = ethrpc.NewWithDefaults(cfg.Endpoint)
	if err != nil {
		return err
	}

	ctr, ok := multicallContract[net]
	if !ok {
		return ErrNoMulticall
	}

	MulticallClient, err = multicall.New(ETHClient, multicall.ContractAddress(ctr))
	if err != nil {
		return err
	}

	return nil
}

// Set sets the clients according to the activeConf
func Set() error {
	cm.Lock()
	defer cm.Unlock()

	config, err := config.GetActiveConf()
	if err != nil {
		return err
	}

	Client, err = ethclient.Dial(config.Endpoint)
	if err != nil {
		return err
	}

	ETHClient, err = ethrpc.NewWithDefaults(config.Endpoint)
	if err != nil {
		return err
	}

	MulticallClient, err = multicall.New(ETHClient, multicall.ContractAddress(multicallContract[config.Name]))
	if err != nil {
		return err
	}

	return nil
}

// LoadAndGet does the same as Load(net string) error but return the pointer from activeConfig
func LoadAndGet(net string) (*config.ActiveConfig, error) {
	err := config.SetActiveConf(net)
	if err != nil {
		return nil, err
	}

	cfg, err := config.GetActiveConf()
	if err != nil {
		return nil, err
	}

	cm.Lock()
	defer cm.Unlock()

	Client, err = ethclient.Dial(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	ETHClient, err = ethrpc.NewWithDefaults(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	ctr, ok := multicallContract[net]
	if !ok {
		return nil, ErrNoMulticall
	}

	MulticallClient, err = multicall.New(ETHClient, multicall.ContractAddress(ctr))
	if err != nil {
		return nil, err
	}

	return cfg, nil

}
