package handler

import (
	"context"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/jon4hz/emergenyWithdrawer/internal/client"
	"github.com/jon4hz/emergenyWithdrawer/internal/config"
)

func SetCurrentView(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	view := v.Name()

	_, err := g.SetCurrentView(view)
	if err != nil {
		return err
	}

	err = cursorCtrl(g, v)
	if err != nil {
		return err
	}

	return nil
}

func cursorCtrl(g *gocui.Gui, v *gocui.View) error {
	var view string
	if v != nil {
		view = v.Name()
	}

	switch view {
	case "endpoint", "gasprice", "gaslimit", "privatekey":
		g.Cursor = true
	default:
		g.Cursor = false
	}
	return nil
}

func SetConf() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v == nil {
			return nil
		}

		view := v.Name()

		msgLogChan <- fmt.Sprintf("set active network to %s\n", view)

		var first bool

		cfg, err := config.GetActiveConf()
		if err != nil {
			if err.Error() != config.ErrNoActiveConf.Error() {
				return err
			}

			err = config.SetActiveConf(view)
			if err != nil {
				errLogChan <- err
				return nil
			}
			cfg, err = config.GetActiveConf()
			if err != nil {
				errLogChan <- err
				return nil
			}
			first = true
		}

		if cfg.Name != view || first {

			prevNet := cfg.Name

			err := config.SetActiveConf(view)
			if err != nil {
				return err
			}
			cfg, err = config.GetActiveConf()
			if err != nil {
				return err
			}

			g.Update(func(g *gocui.Gui) error {
				v, err := g.View(prevNet)
				if err != nil {
					return err
				}
				v.FrameRunes = []rune{}

				v, err = g.View(view)
				if err != nil {
					return err
				}

				v.FrameRunes = []rune{'═', '║', '╔', '╗', '╚', '╝'}

				return nil
			})

			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("gaslimit")
				if err != nil {
					return err
				}

				v.Clear()
				fmt.Fprint(v, cfg.GasLimit)

				v, err = g.View("endpoint")
				if err != nil {
					return err
				}

				v.Clear()
				fmt.Fprint(v, cfg.Endpoint)

				return nil
			})

			err = client.Load(view)
			if err != nil {
				return err
			}

			go func(g *gocui.Gui, v *gocui.View) {
				gas, err := client.Client.SuggestGasPrice(context.TODO())
				if err != nil {
					return
				}
				g.Update(func(g *gocui.Gui) error {
					v, err := g.View("gasprice")
					if err != nil {
						return err
					}

					v.Clear()
					fmt.Fprint(v, gas)

					return nil
				})
			}(g, v)

			return nil
		} else if cfg.Name == view {
			go func(g *gocui.Gui, v *gocui.View) {
				gas, err := client.Client.SuggestGasPrice(context.TODO())
				if err != nil {
					return
				}
				g.Update(func(g *gocui.Gui) error {
					v, err := g.View("gasprice")
					if err != nil {
						return err
					}

					v.Clear()
					fmt.Fprint(v, gas)

					return nil
				})
			}(g, v)
		}

		return nil
	}
}

func GetPools() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}
}
