package client

import (
	"errors"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jon4hz/emergenyWithdrawer/internal/config"
	"github.com/jon4hz/web3-multicall-go/multicall"
)

var (
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
