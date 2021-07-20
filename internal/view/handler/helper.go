package handler

import "github.com/awesome-gocui/gocui"

func cursorCtrl(g *gocui.Gui, v *gocui.View) error {
	var view string
	if v != nil {
		view = v.Name()
	}

	switch view {
	case "endpoint", "gasprice", "gaslimit", "privatekey", "masterchef", "pool", "log":
		g.Cursor = true
	default:
		g.Cursor = false
	}
	return nil
}

func getSelectedFrameRunes() []rune {
	return []rune{'═', '║', '╔', '╗', '╚', '╝'}
}

func getUnselectedFrameRunes() []rune {
	return []rune{}
}

// markViewValid marks the view as valid with a green frame
func markViewValid(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v.FrameColor = gocui.ColorGreen
		return nil
	})
	return nil
}

// markViewInvalid marks the view as invalid with a red frame
func markViewInvalid(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		v.FrameColor = gocui.ColorRed
		return nil
	})
	return nil
}
