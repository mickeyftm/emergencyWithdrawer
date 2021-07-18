package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/jon4hz/emergenyWithdrawer/internal/client"
	"github.com/jon4hz/emergenyWithdrawer/internal/config"
	"github.com/jon4hz/emergenyWithdrawer/internal/util"
)

func InputEndpoint() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vbuf := v.ViewBuffer()

		if !util.URLMatch.MatchString(vbuf) {

			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("endpoint")
				if err != nil {
					return err
				}
				v.FrameColor = gocui.ColorRed
				return nil
			})
			return nil
		}

		g.Update(func(g *gocui.Gui) error {
			v.FrameColor = gocui.ColorGreen
			return nil
		})

		err := updateEndpoint(g, vbuf)
		if err != nil {
			return err
		}

		return nil
	}
}

func updateEndpoint(g *gocui.Gui, endpoint string) error {
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

	msgLogChan <- "set endpoint to " + endpoint

	return nil
}

func updateGasPrice(g *gocui.Gui, v *gocui.View) error {
	var failed bool
	gasPrice, err := client.Client.SuggestGasPrice(context.TODO())
	if err != nil {
		errLogChan <- errors.New("failed to get gas price: " + err.Error())
		failed = true
	}

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("info")
		if err != nil {
			return err
		}
		if !failed {
			price, _ := util.Gwei2Eth(gasPrice).Float64()
			v.SetLine(6, fmt.Sprintf("\tGas Price: %.18f\n", price))

			return nil
		}

		v.SetLine(6, fmt.Sprintf("\tGas Price: %s\n", "failed to get gas price"))

		return nil
	})

	return nil
}

func InputGasPrice() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vbuf := v.ViewBuffer()
		_ = vbuf
		return nil
	}
}

func InputGasLimit() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vbuf := v.ViewBuffer()
		_ = vbuf
		return nil
	}
}

func InputPrivateKey() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v == nil {
			return nil
		}
		vbuf := v.ViewBuffer()
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
		return nil
	}
}

func InputMasterchef() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vbuf := v.ViewBuffer()
		_ = vbuf
		return nil
	}
}
