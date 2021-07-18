package handler

import "github.com/awesome-gocui/gocui"

func cursorCtrl(g *gocui.Gui, v *gocui.View) error {
	var view string
	if v != nil {
		view = v.Name()
	}

	switch view {
	case "endpoint", "gasprice", "gaslimit", "privatekey", "masterchef", "pool":
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
