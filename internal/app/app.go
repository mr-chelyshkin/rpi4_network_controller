package app

import (
	"context"
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/ui"
	"os"

	"github.com/gdamore/tcell/v2"
)

type cmdConnectNetworkDetails struct {
	form     func()
	subTitle string
	title    string
}

// Run application.
func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//
	values := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	ctx = context.WithValue(ctx, rpi4_network_controller.CtxKeyHotkeys, values)

	stop := make(chan struct{}, 1)
	defer close(stop)

	data := [][]string{
		{"Connect", "connect to wifi network"},
		{"Disconnect", "interrupt wifi connection"},
	}
	frame := ui.ContentTable(ctx, data)

	//frame := tview.NewList().
	//	AddItem("Connect", "connect to wifi network", '1', func() { cmdConnect(stop) }).
	//	AddItem("Disconnect", "interrupt wifi connection", '2', func() {})

	ui.App.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				stop <- struct{}{}
				Run()
			case tcell.KeyCtrlC:
				stop <- struct{}{}
				ui.App.Stop()
				os.Exit(0)
			}
			return event
		},
	)
	return ui.Start(ctx, frame)
}
