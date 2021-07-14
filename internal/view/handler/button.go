package handler

import (
	"github.com/awesome-gocui/gocui"
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

		cfg, err := config.GetActiveConf()
		if err != nil {
			if err.Error() != config.ErrNoActiveConf.Error() {
				return err
			}

			err = config.SetActiveConf(view)
			if err != nil {
				// TODO replace return err with write err to log, if return err hits here the code panics
				return err
			}
			cfg, err = config.GetActiveConf()
			if err != nil {
				// TODO replace return err with write err to log, if return err hits here the code panics
				return err
			}
		}

		if cfg.Name != view {
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View(cfg.Name)
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

			err := config.SetActiveConf(view)
			if err != nil {
				return err
			}

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
