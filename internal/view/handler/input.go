package handler

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jon4hz/emergenyWithdrawer/internal/client"
	"github.com/jon4hz/emergenyWithdrawer/internal/config"
	"github.com/jon4hz/emergenyWithdrawer/internal/logging"
	"github.com/jon4hz/emergenyWithdrawer/internal/secrets"
	"github.com/jon4hz/emergenyWithdrawer/internal/utils"
)

func InputEndpoint() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vbuf := v.ViewBuffer()

		if !utils.URLMatch.MatchString(vbuf) {

			markViewInvalid(g, v)
			return nil
		}

		markViewValid(g, v)
		setEndpoint(g, vbuf)

		logging.InfoChan <- "set endpoint to " + vbuf

		return nil
	}
}

func setEndpoint(g *gocui.Gui, endpoint string) error {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("info")
		if err != nil {
			return err
		}
		v.SetLine(5, fmt.Sprintf("\tEndpoint: %s\n", endpoint))
		return nil
	})

	err := config.SetActiveEndpoint(endpoint)
	if err != nil {
		return err
	}

	err = client.Set()
	if err != nil {
		return err
	}

	err = updateGasPrice(g, nil)
	if err != nil {
		return err
	}

	return nil
}

func updateEndpoint(g *gocui.Gui, endpoint string) error {

	err := setEndpoint(g, endpoint)
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("endpoint")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprint(v, endpoint)

		return nil
	})

	return nil
}

func InputGasPrice() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vbuf := v.ViewBuffer()

		gasPrice, err := strconv.ParseFloat(vbuf, 64)
		if err != nil {
			logging.ErrChan <- errors.New("failed to parse gas price: " + err.Error())
			markViewInvalid(g, v)
			return nil
		}

		price := utils.FloatToBigInt(gasPrice)

		markViewValid(g, v)
		setGasPrice(g, price)

		logging.InfoChan <- fmt.Sprintf("set gas price to %.18f", gasPrice)

		return nil
	}
}

func setGasPrice(g *gocui.Gui, gasPrice *big.Int) error {

	config.SetActiveGasPrice(gasPrice)

	price, _ := utils.Gwei2Eth(gasPrice).Float64()
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("info")
		if err != nil {
			return err
		}
		v.SetLine(6, fmt.Sprintf("\tGas Price: %.18f\n", price))
		return nil
	})

	return nil
}

func updateGasPrice(g *gocui.Gui, v *gocui.View) error {
	var failed bool

	if client.Client == nil {
		logging.ErrChan <- errors.New("client is nil, did you set an endpoint?")
		return nil
	}

	gasPrice, err := client.Client.SuggestGasPrice(context.TODO())
	if err != nil {
		logging.ErrChan <- errors.New("failed to get gas price: " + err.Error())
		failed = true
	}

	if !failed {
		err = setGasPrice(g, gasPrice)
		if err != nil {
			return err
		}
	}

	g.Update(func(g *gocui.Gui) error {
		if !failed {
			v, err := g.View("gasprice")
			if err != nil {
				return err
			}
			v.Clear()
			price, _ := utils.Gwei2Eth(gasPrice).Float64()
			fmt.Fprintf(v, "%.18f", price)

			return nil
		}

		v.SetLine(6, fmt.Sprintf("\tGas Price: %s\n", "failed to get suggested gas price"))

		logging.ErrChan <- errors.New("failed to get suggested gas price")

		return nil
	})

	return nil
}

func InputGasLimit() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vbuf := v.ViewBuffer()

		gasLimit, err := strconv.ParseUint(vbuf, 10, 64)
		if err != nil {
			logging.ErrChan <- errors.New("failed to parse gas limit: " + err.Error())
			markViewInvalid(g, v)
			return nil
		}

		markViewValid(g, v)
		setGasLimit(g, gasLimit)

		logging.InfoChan <- "set gas limit to " + vbuf

		return nil
	}
}

func setGasLimit(g *gocui.Gui, gasLimit uint64) error {

	config.SetActiveGasLimit(gasLimit)

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("info")
		if err != nil {
			return err
		}
		v.SetLine(7, fmt.Sprintf("\tGas Limit: %d\n", gasLimit))
		return nil
	})

	return nil
}

func updateGasLimit(g *gocui.Gui, gasLimit uint64) error {

	err := setGasLimit(g, gasLimit)
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("gaslimit")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintf(v, "%d", gasLimit)

		return nil
	})
	return nil
}

func InputPrivateKey() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {

		vbuf := v.ViewBuffer()

		err := setPrivateKey(g, vbuf)
		if err != nil {
			logging.ErrChan <- errors.New("failed to set private key: " + err.Error())
			markViewInvalid(g, v)
			return nil
		}

		markViewValid(g, v)

		// replace input with asteriks
		lenKey := len(vbuf)
		var astrisks string
		for i := 0; i < lenKey; i++ {
			astrisks += "*"
		}
		g.Update(func(g *gocui.Gui) error {
			cx, cy := v.Cursor()
			v.Clear()
			fmt.Fprint(v, astrisks)
			err := v.SetCursor(cx, cy)
			if err != nil {
				return err
			}
			return nil
		})

		key, err := secrets.GetPrivateKey()
		if err != nil {
			return err
		}

		logging.InfoChan <- "set private key for address " + utils.PrivateKeyToAddress(key)

		return nil
	}
}

func setPrivateKey(g *gocui.Gui, key string) error {

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return err
	}

	err = secrets.SetPrivateKeyString(key)
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("info")
		if err != nil {
			return err
		}
		v.SetLine(1, fmt.Sprintf("\tAddress: %s\n", utils.PrivateKeyToAddress(privateKey)))
		return nil
	})
	return nil
}

func InputMasterchef() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vbuf := v.ViewBuffer()
		_ = vbuf
		return nil
	}
}
