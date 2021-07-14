package view

import "github.com/awesome-gocui/gocui"

type PoolWidget struct {
	name    string
	getView func(*gocui.Gui) (int, int, int, int)
}

type ButtonWidget struct {
	name    string
	label   string
	getView func(*gocui.Gui) (int, int, int, int)
	opts    *ButtonWidgetOpts
	handler func(g *gocui.Gui, v *gocui.View) error
}

type ButtonWidgetOpts struct {
	frameColor   gocui.Attribute
	fgColor      gocui.Attribute
	frameRunes   []rune
	textCentered bool
}

type InputWidget struct {
	name    string
	title   string
	getView func(*gocui.Gui) (int, int, int, int)
	opts    *InputWidgetOpts
	handler func(g *gocui.Gui, v *gocui.View) error
}
type InputWidgetOpts struct {
	frameColor gocui.Attribute
	fgColor    gocui.Attribute
}
