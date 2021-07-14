package handler

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/gdamore/tcell/v2"
	"github.com/jon4hz/emergenyWithdrawer/internal/helper/regex"
)

func InputEndpoint() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vbuf := v.ViewBuffer()

		if !regex.URLMatch.MatchString(vbuf) {

			fmt.Println("no match")
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("endpoint")
				if err != nil {
					return err
				}
				v.FrameColor = gocui.Attribute(tcell.ColorRed)
				return nil
			})
		}
		g.Update(func(g *gocui.Gui) error {
			v.FrameColor = gocui.ColorGreen
			return nil
		})
		return nil
	}
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
