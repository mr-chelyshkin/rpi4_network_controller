package app

import (
	"context"

	"github.com/mr-chelyshkin/rpi4_network_controller"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Run() error {
	stop := make(chan struct{}, 1)
	ctx := context.Background()

	main := tview.NewList().
		AddItem("Connect", "connect to wifi network", '1', func() { cmdConnect(ctx, stop) }).
		AddItem("Disconnect", "interrupt wifi connection", '2', nil)
	frameMain := frameDefault(ctx, main, nil)

	rpi4_network_controller.App.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				rpi4_network_controller.App.SetRoot(frameMain, true)
				stop <- struct{}{}
			case tcell.KeyCtrlC:
				rpi4_network_controller.App.Stop()
			}
			return event
		},
	)
	return rpi4_network_controller.App.SetRoot(
		frameMain, true,
	).SetFocus(
		frameMain,
	).Run()
}
