package handler

import "github.com/awesome-gocui/gocui"

func SetCurrentView(g *gocui.Gui, v *gocui.View) error {
	var view string
	if v != nil {
		view = v.Name()
	}
	_, err := g.SetCurrentView(view)
	if err != nil {
		return err
	}
	err = cursorCtrl(g, v)
	return err
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
		return nil
	}
}

func GetPools() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}
}
