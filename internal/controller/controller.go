package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Run() error {
	app := tview.NewApplication()

	main := tview.NewList().
		AddItem("Connect", "connect to wifi network", '1', func() { cmdConnect(app) }).
		AddItem("Disconnect", "interrupt wifi connection", '2', nil)
	frameMain := frameDefault(main)

	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				app.SetRoot(frameMain, true)
			case tcell.KeyCtrlC:
				app.Stop()
			}
			return event
		},
	)
	return app.SetRoot(frameMain, true).SetFocus(frameMain).Run()
}
