package handler

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/jon4hz/emergenyWithdrawer/internal/client"
	"github.com/jon4hz/emergenyWithdrawer/internal/config"
)

func SetCurrentView(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	var prevView string
	prevV := g.CurrentView()
	if prevV != nil {
		prevView = prevV.Name()
	}

	view := v.Name()
	if prevView != view {

		if prevV != nil {
			prevV.TitleColor = gocui.ColorDefault
		}

		v.TitleColor = gocui.ColorCyan
		_, err := g.SetCurrentView(view)
		if err != nil {
			return err
		}

		err = cursorCtrl(g, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetConf() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if v == nil {
			return nil
		}

		err := SetCurrentView(g, v)
		if err != nil {
			errLogChan <- err
			return nil
		}

		view := v.Name()

		var first bool
		cfg, err := config.GetActiveConf()
		if err != nil {
			if err.Error() != config.ErrNoActiveConf.Error() {
				return err
			}
			first = true
		}

		var name string
		if cfg != nil {
			name = cfg.Name
		}

		if name != view || first {

			if name != "" {
				markButtonUnselected(g, name)
			}

			markButtonSelected(g, view)

			cfg, err := client.LoadAndGet(view)
			if err != nil {
				return err
			}

			err = updateGasPrice(g, v)
			if err != nil {
				return err
			}

			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("gaslimit")
				if err != nil {
					return err
				}

				v.Clear()
				fmt.Fprintf(v, "%d", cfg.GasLimit)

				v, err = g.View("endpoint")
				if err != nil {
					return err
				}

				v.Clear()
				fmt.Fprint(v, cfg.Endpoint)

				return nil
			})

			err = updateEndpoint(g, cfg.Endpoint)
			if err != nil {
				return err
			}

			msgLogChan <- fmt.Sprintf("set active network to %s", view)

			return nil

		}

		return nil
	}
}

func GetPools() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}
}

func Withdraw() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}
}

func markButtonSelected(g *gocui.Gui, button string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(button)
		if err != nil {
			return err
		}
		v.FrameRunes = getSelectedFrameRunes()
		return nil
	})
}

func markButtonUnselected(g *gocui.Gui, button string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(button)
		if err != nil {
			return err
		}
		v.FrameRunes = getUnselectedFrameRunes()
		return nil
	})
}
