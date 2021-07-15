package secrets

import (
	"crypto/ecdsa"
)

var (
	ErrSecretNotFound = "private key not found in keyring"

	PrivateKey *ecdsa.PrivateKey
)

/* func init() {

}

func Test() {
	service := "emercenyWithdrawer_privateKey"
	user, err := user.Current()
	if err != nil {
		log.Panic(err)
	}

	username := user.Username */
/*
	password := "test1"

	 err = keyring.Set(service, username, password)
	if err != nil {
		log.Fatal(err)
	} */

/* 	var configurePrivateKey bool
	secret, err := keyring.Get(service, username)
	if err != nil {
		if err.Error() != ErrSecretNotFound {
			log.Fatal(err)
		} else {

			fmt.Print("No private key found. Would you like to set one? [Y/n] ")
			err = nil
			configurePrivateKey = getValidConfirmation()

		}
	}

	if configurePrivateKey {

		fmt.Print("Please enter your private key: ")

		PrivateKey = getPrivateKey()
	}

	_ = secret
} */
/*
func getValidConfirmation() bool {
	var answer string
	invalid := true

	for invalid {
		fmt.Scanln(&answer)

		confirm, err := helper.GetConfirmation(answer)
		if err != nil {
			fmt.Print("invalid input, try again: ")
			continue
		}

		return confirm

	}
	return false
}

func getPrivateKey() *ecdsa.PrivateKey {
	var answer string
	invalid := true

	for invalid {
		fmt.Scanln(&answer)

		privateKey, err := crypto.HexToECDSA(answer)
		if err != nil {
			fmt.Print("invalid input, try again: ")
			continue
		}

		return privateKey
	}
	return nil
}
*/
