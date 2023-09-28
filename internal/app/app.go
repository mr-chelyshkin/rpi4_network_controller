package app

import (
	"context"

	"github.com/gdamore/tcell/v2"
)

func Run() error {
	frame := frameMain(context.Background())
	stop := make(chan struct{}, 1)

	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				app.SetRoot(frame, true)
				stop <- struct{}{}
			case tcell.KeyCtrlC:
				app.Stop()
			}
			return event
		},
	)
	return app.SetRoot(frame, true).SetFocus(frame).Run()
}
