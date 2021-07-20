package utils

import "math/big"

func Gwei2Eth(x *big.Int) *big.Float {
	po := new(big.Float).SetInt(
		new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
	)

	y := new(big.Float).SetInt(x)
	return new(big.Float).Quo(y, po)
}
