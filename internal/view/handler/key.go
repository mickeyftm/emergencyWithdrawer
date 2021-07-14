package handler

import "github.com/awesome-gocui/gocui"

var (
	DonePoolFeed = make(chan struct{})
)

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		/* cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		} */
		v.MoveCursor(0, 1)
	}
	return nil
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		/* ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		} */
		v.MoveCursor(0, -1)
	}
	return nil
}

func Quit(g *gocui.Gui, v *gocui.View) error {
	close(DonePoolFeed)
	return gocui.ErrQuit
}
