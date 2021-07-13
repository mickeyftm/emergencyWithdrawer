package view

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/jon4hz/emergenyWithdrawer/internal/view/handler"
)

var (
	donePoolFeed = make(chan struct{})
	wg           sync.WaitGroup
)

func New() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Mouse = true

	setLayout(g)

	if err := setKeybindings(g); err != nil {
		log.Panicln(err)
	}

	wg.Add(1)
	go updatePools(g)

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

	wg.Wait()
}

func setLayout(g *gocui.Gui) {

	ethButton := NewButtonWidget(
		"eth",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return 0, 0, (maxX / 18) - 1, 2
		},
		"ETH",
		&ButtonWidgetOpts{
			FrameColor: gocui.ColorBlue,
		},
		handler.SetConf(),
	)
	bscButton := NewButtonWidget(
		"bsc",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return (maxX / 18) + 1, 0, (maxX / 9) - 1, 2
		},
		"BSC",
		&ButtonWidgetOpts{
			FrameColor: gocui.ColorYellow,
		},
		handler.SetConf(),
	)
	maticButton := NewButtonWidget(
		"matic",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return (maxX / 9) + 1, 0, 3*(maxX/18) + 1, 2
		},
		"MATIC",
		&ButtonWidgetOpts{
			FrameColor: gocui.GetColor("#8A2BE2"),
		},
		handler.SetConf(),
	)

	searchPoolsButton := NewButtonWidget(
		"searchPools",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return maxX / 3, 0, 2*(maxX/3) - 1, 2
		},
		"Search for pools",
		&ButtonWidgetOpts{
			FrameColor: gocui.ColorBlue,
		},
		getPools(),
	)
	pools := NewPoolWidget(
		"pools",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, maxY := g.Size()
			return maxX / 3, 3, 2*(maxX/3) - 1, maxY/3 - 1
		},
	)

	g.SetManager(ethButton, bscButton, maticButton, searchPoolsButton, pools)
}

func setKeybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		return err
	}
	/* if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, toggleButton); err != nil {
		return err
	} */
	err = g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, handler.SetCurrentView)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("pools", gocui.KeyArrowUp, gocui.ModNone, handler.CursorUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("pools", gocui.KeyArrowDown, gocui.ModNone, handler.CursorDown)
	if err != nil {
		return err
	}

	return nil
}

func NewPoolWidget(name string, getView func(*gocui.Gui) (int, int, int, int)) *PoolWidget {
	return &PoolWidget{name: name, getView: getView}
}

func (w *PoolWidget) Layout(g *gocui.Gui) error {
	x0, y0, x1, y1 := w.getView(g)
	v, err := g.SetView(w.name, x0, y0, x1, y1, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "Pools"
		v.Autoscroll = true
		v.Highlight = true
		v.SelFgColor = gocui.ColorYellow
		v.Wrap = true
	}
	return nil
}

func NewButtonWidget(name string, getView func(*gocui.Gui) (int, int, int, int), label string, opts *ButtonWidgetOpts, handler func(g *gocui.Gui, v *gocui.View) error) *ButtonWidget {

	if opts == nil {
		opts = new(ButtonWidgetOpts)
	}

	return &ButtonWidget{
		name:    name,
		getView: getView,
		label:   label,
		opts:    opts,
		handler: handler,
	}
}

func (w *ButtonWidget) Layout(g *gocui.Gui) error {
	x0, y0, x1, y1 := w.getView(g)

	v, err := g.SetView(w.name, x0, y0, x1, y1, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.FrameColor = w.opts.FrameColor
		v.FgColor = w.opts.FgColor
		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
			return err
		}
		text := w.label
		if w.opts.TextCentered {
			text = centerText(w.label, x1-x0)
		}

		fmt.Fprint(v, text)
	}
	return nil
}

func centerText(text string, lenght int) string {
	if len(text) >= lenght {
		return text
	}

	textRune := []rune(text)
	whitespaces := (lenght - len(text)) / 2
	whitespace := " "
	for i := 0; i < whitespaces-1; i++ {
		textRune = append([]rune(whitespace), textRune...)
	}

	return string(textRune)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(donePoolFeed)
	return gocui.ErrQuit
}

/* func toggleButton(g *gocui.Gui, v *gocui.View) error {
	nextview := "butdown"
	if v != nil && v.Name() == "butdown" {
		nextview = "butup"
	}
	_, err := g.SetCurrentView(nextview)
	return err
} */

func getPools() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}
}

func updatePools(g *gocui.Gui) {
	defer wg.Done()
	ticker := time.NewTicker(time.Second * 1)
	i := 0
	for {
		select {
		case <-donePoolFeed:
			return
		case <-ticker.C:
			/* out, err := g.View("pools")
			if err != nil {
				continue
			}
			fmt.Fprintln(out, <-poolFeed) */
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("pools")
				if err != nil {
					return err
				}
				fmt.Fprintln(v, i)
				return nil
			})
			i++

			if i == 10 {
				ticker.Stop()
			}
		}

	}
}
