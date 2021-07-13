package handler

import "github.com/awesome-gocui/gocui"

func SetCurrentView(g *gocui.Gui, v *gocui.View) error {
	var view string
	if v != nil {
		view = v.Name()
	}
	_, err := g.SetCurrentView(view)
	return err
}

func SetConf() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}
}
