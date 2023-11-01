package app

import (
	"context"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

	stop := make(chan struct{}, 1)
	defer close(stop)

	frame := tview.NewList().
		AddItem("Connect", "connect to wifi network", '1', func() { cmdConnect(stop) }).
		AddItem("Disconnect", "interrupt wifi connection", '2', func() {})

	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				stop <- struct{}{}
				Run()
			case tcell.KeyCtrlC:
				stop <- struct{}{}
				app.Stop()
				os.Exit(0)
			}
			return event
		},
	)
	return appRun(frameWrapper(ctx, frame, nil))
}
