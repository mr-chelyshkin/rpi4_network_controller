package internal

import (
	"context"

	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/schedule"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/ui"
	"github.com/rivo/tview"
)

func connect(ctx context.Context, interrupt chan struct{}) {
	output := make(chan string, 1)
	defer close(output)

	networks := make(chan []map[string]string)
	defer close(networks)

	ctx = context.WithValue(ctx, rpi4_network_controller.CtxKeyOutputCh, output)
	ctx, cancel := context.WithCancel(ctx)

	view := tview.NewList()
	go ui.Draw(ctx, view)
	go func() {
		for {
			select {
			case networks := <-networks:
				ui.App.QueueUpdateDraw(func() {
					view.Clear()

					for _, network := range networks {
						network := network
						view.AddItem(network["ssid"], network["level"], '*',
							func() {

							})
					}
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	schedule.NetworkScan(ctx, networks)

	func() {
		select {
		case <-ctx.Done():
			return
		case <-interrupt:
			cancel()
			return
		}
	}()
}
