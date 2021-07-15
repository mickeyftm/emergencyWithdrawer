package handler

import (
	"github.com/awesome-gocui/gocui"
)

var (
	DonePoolFeed = make(chan struct{})
)

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	close(DonePoolFeed)
	return gocui.ErrQuit
}

func ToggleInput(g *gocui.Gui, v *gocui.View) error {
	var nextview string
	if v != nil {
		switch v.Name() {
		case "endpoint":
			nextview = "gasprice"
		case "gasprice":
			nextview = "gaslimit"
		case "gaslimit":
			nextview = "privatekey"
		case "privatekey":
			nextview = "masterchef"
		case "masterchef":
			nextview = "pool"
		default:
			nextview = "endpoint"
		}
	}
	/* _, err := g.SetCurrentView(nextview)
	return err */

	v, err := g.View(nextview)
	if err != nil {
		return err
	}

	err = SetCurrentView(g, v)
	if err != nil {
		return err
	}
	return nil
}
