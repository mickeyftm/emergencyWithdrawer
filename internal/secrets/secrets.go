package secrets

import (
	"crypto/ecdsa"
	"errors"
	"os/user"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zalando/go-keyring"
)

var (
	ErrSecretNotFound = errors.New("private key not found in keyring")
	ErrSecretIsNil    = errors.New("private key is nil")

	privateKey *ecdsa.PrivateKey

	username string
)

const (
	service = "emergencyWithdrawer_privateKey"
)

func init() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	username = user.Username

}

func GetPrivateKey() (*ecdsa.PrivateKey, error) {
	if privateKey == nil {
		var err error
		privateKey, err = getPrivateKeyFromKeyring()
		if err != nil {
			return nil, err
		}
	}
	return privateKey, nil
}

func getPrivateKeyFromKeyring() (*ecdsa.PrivateKey, error) {
	secret, err := keyring.Get(service, username)
	if err != nil {
		return nil, err
	}
	key, err := crypto.HexToECDSA(secret)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func SetPrivateKeyString(key string) error {
	var err error
	privateKey, err = crypto.HexToECDSA(key)
	if err != nil {
		return err
	}

	err = storePrivateKeyInKeyring(key)
	if err != nil {
		return err
	}

	return nil
}

func storePrivateKeyInKeyring(key string) error {
	return keyring.Set(service, username, key)
}
