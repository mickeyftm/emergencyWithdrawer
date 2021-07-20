package utils

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address
}
