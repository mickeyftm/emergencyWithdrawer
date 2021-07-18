package handler

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

var (
	msgLogChan  = make(chan string)
	warnLogChan = make(chan string)
	errLogChan  = make(chan error)
)

func FeedLog(g *gocui.Gui, msgLog <-chan string, warnLog <-chan string, errLog <-chan error) {
	for {
		select {
		case log := <-msgLog:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("log")
				if err != nil {
					return err
				}
				fmt.Fprintln(v, "[info]", log)
				return nil
			})
		case log := <-warnLog:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("log")
				if err != nil {
					return err
				}
				fmt.Fprintln(v, "[warn]", log)
				return nil
			})
		case log := <-errLog:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("log")
				if err != nil {
					return err
				}
				fmt.Fprintln(v, "[erro]", log)
				return nil
			})
		}

	}
}

// GetMsgLogChan returns the channel used to send messages to the log view
func GetMsgLogChan() chan string {
	return msgLogChan
}

// GetWarnLogChan returns the channel used to send warnings to the log view
func GetWarnLogChan() chan string {
	return warnLogChan
}

// GetErrLogChan returns the channel used to send errors to the log view
func GetErrLogChan() chan error {
	return errLogChan
}
