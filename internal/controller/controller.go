package controller

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Run() error {
	stop := make(chan struct{}, 1)
	app := tview.NewApplication()
	ctx := context.Background()

	main := tview.NewList().
		AddItem("Connect", "connect to wifi network", '1', func() { cmdConnect(ctx, stop, app) }).
		AddItem("Disconnect", "interrupt wifi connection", '2', nil)
	frameMain := frameDefault(ctx, main, nil)

	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				app.SetRoot(frameMain, true)
				stop <- struct{}{}
			case tcell.KeyCtrlC:
				app.Stop()
			}
			return event
		},
	)
	return app.SetRoot(frameMain, true).SetFocus(frameMain).Run()
}
