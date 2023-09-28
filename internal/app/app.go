package app

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
				app.SetRoot(frameWrapper(ctx, frame, nil), true).SetFocus(frame)
				stop <- struct{}{}
			case tcell.KeyCtrlC:
				app.Stop()
			}
			return event
		},
	)
	return app.SetRoot(frameWrapper(ctx, frame, nil), true).SetFocus(frame).Run()
}
