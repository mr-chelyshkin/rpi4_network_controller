package app

import (
	"context"

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
	stop := make(chan struct{}, 1)
	ctx := context.Background()

	frame := tview.NewList().
		AddItem("Connect", "connect to wifi network", '1', func() { cmdConnect(stop) }).
		AddItem("Disconnect", "interrupt wifi connection", '2', func() {})

	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				setFrame(frameWrapper(ctx, frame, nil))
				stop <- struct{}{}
			case tcell.KeyCtrlC:
				stop <- struct{}{}
				app.Stop()
			}
			return event
		},
	)
	return appRun(frameWrapper(ctx, frame, nil))
}
