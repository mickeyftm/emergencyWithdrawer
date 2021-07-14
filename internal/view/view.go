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
	wg sync.WaitGroup
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
			frameColor: gocui.ColorBlue,
		},
		handler.SetConf(),
	)

	bscButton := NewButtonWidget(
		"bsc",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return (maxX / 18), 0, 2*(maxX/18) - 1, 2
		},
		"BSC",
		&ButtonWidgetOpts{
			frameColor: gocui.ColorYellow,
		},
		handler.SetConf(),
	)

	maticButton := NewButtonWidget(
		"matic",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return 2 * (maxX / 18), 0, 3*(maxX/18) - 1, 2
		},
		"MATIC",
		&ButtonWidgetOpts{
			frameColor: gocui.ColorCyan,
		},
		handler.SetConf(),
	)

	ftmButton := NewButtonWidget(
		"ftm",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return 3 * (maxX / 18), 0, 4*(maxX/18) - 1, 2
		},
		"FTM",
		&ButtonWidgetOpts{
			frameColor: gocui.ColorBlue,
		},
		handler.SetConf(),
	)

	kccButton := NewButtonWidget(
		"kcc",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return 4 * (maxX / 18), 0, 5*(maxX/18) - 1, 2
		},
		"KCC",
		&ButtonWidgetOpts{
			frameColor: gocui.ColorGreen,
		},
		handler.SetConf(),
	)

	rpcInput := NewInputWidget(
		"endpoint",
		"Endpoint",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return 0, 3, maxX/3 - 1 - 1, 5
		},
		nil,
		handler.InputEndpoint(),
	)

	gasPriceInput := NewInputWidget(
		"gasprice",
		"Gas Price",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return 0, 6, maxX/6 - 2, 8
		},
		nil,
		handler.InputGasPrice(),
	)

	gasLimitInput := NewInputWidget(
		"gaslimit",
		"Gas Limit",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return maxX / 6, 6, maxX/3 - 2, 8
		},
		nil,
		handler.InputGasPrice(),
	)

	privateKeyInput := NewInputWidget(
		"privatekey",
		"Private Key",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return 0, 9, maxX/3 - 2, 11
		},
		nil,
		handler.InputPrivateKey(),
	)

	searchPoolsButton := NewButtonWidget(
		"searchPools",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, _ := g.Size()
			return maxX / 3, 0, 2*(maxX/3) - 1, 2
		},
		"Search for pools",
		&ButtonWidgetOpts{
			frameColor: gocui.ColorBlue,
		},
		handler.GetPools(),
	)

	pools := NewPoolWidget(
		"pools",
		func(*gocui.Gui) (int, int, int, int) {
			maxX, maxY := g.Size()
			y1 := maxY/3 - 1
			if y1 < 11 {
				y1 = 11
			}
			return maxX / 3, 3, 2*(maxX/3) - 1, y1
		},
	)

	g.SetManager(
		ethButton,
		bscButton,
		maticButton,
		ftmButton,
		kccButton,
		rpcInput,
		gasPriceInput,
		gasLimitInput,
		privateKeyInput,
		searchPoolsButton,
		pools,
	)
}

func setKeybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, handler.Quit)
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

		if len(w.opts.frameRunes) > 0 {
			v.FrameRunes = w.opts.frameRunes
		}

		v.FrameColor = w.opts.frameColor
		v.FgColor = w.opts.fgColor

		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.handler); err != nil {
			return err
		}
		text := w.label
		if w.opts.textCentered {
			text = centerText(w.label, x1-x0)
		}

		fmt.Fprint(v, text)
	}
	return nil
}

func NewInputWidget(name, title string, getView func(*gocui.Gui) (int, int, int, int), opts *InputWidgetOpts, handler func(g *gocui.Gui, v *gocui.View) error) *InputWidget {

	if opts == nil {
		opts = &InputWidgetOpts{}
	}

	return &InputWidget{
		name:    name,
		title:   title,
		getView: getView,
		opts:    opts,
		handler: handler,
	}
}

func (w *InputWidget) Layout(g *gocui.Gui) error {
	x0, y0, x1, y1 := w.getView(g)

	v, err := g.SetView(w.name, x0, y0, x1, y1, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = w.title
		v.Editable = true

		v.FrameColor = w.opts.frameColor
		v.FgColor = w.opts.fgColor

		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.handler); err != nil {
			return err
		}

	}
	return nil
}

/* func toggleButton(g *gocui.Gui, v *gocui.View) error {
	nextview := "butdown"
	if v != nil && v.Name() == "butdown" {
		nextview = "butup"
	}
	_, err := g.SetCurrentView(nextview)
	return err
} */

func updatePools(g *gocui.Gui) {
	defer wg.Done()
	ticker := time.NewTicker(time.Second * 1)
	i := 0
	for {
		select {
		case <-handler.DonePoolFeed:
			return
		case <-ticker.C:
			/* g.Update(func(g *gocui.Gui) error {
				v, err := g.View("pools")
				if err != nil {
					return err
				}
				fmt.Fprintln(v, i)
				return nil
			}) */
			i++

			if i == 10 {
				ticker.Stop()
			}
		}

	}
}
