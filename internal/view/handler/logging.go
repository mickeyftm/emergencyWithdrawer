package handler

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

var (
	msgLogChan = make(chan string)
	errLogChan = make(chan error)
)

func FeedLog(g *gocui.Gui) {
	for {
		select {
		case log := <-msgLogChan:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("log")
				if err != nil {
					return err
				}
				fmt.Fprint(v, "[info] ", log)
				return nil
			})
		case log := <-errLogChan:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("log")
				if err != nil {
					return err
				}
				fmt.Fprint(v, "[erro] ", log)
				return nil
			})
		}

	}
}
