package controller

import (
	"context"
	"time"

	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/wifi"

	"github.com/rivo/tview"
)

func cmdConnect(ctx context.Context, cancel context.CancelFunc, app *tview.Application) {
	networks := tview.NewList()
	app.SetRoot(frameDefault(networks), true)

	wifiController, err := wifi.NewWifi()
	if err != nil {
	}
	exec(ctx, cancel, wifiController, app, networks)

	go func(ctx context.Context) {
		ticker := time.NewTicker(rpi4_network_controller.ScanTimeoutSec * time.Second)
		for {
			select {
			case <-ticker.C:
				app.QueueUpdateDraw(
					func() {
						networks.Clear()
						exec(ctx, cancel, wifiController, app, networks)
					},
				)
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

func exec(ctx context.Context, cancel context.CancelFunc, wifi wifi.Wifi, app *tview.Application, list *tview.List) {
	scanResult := wifi.Scan()

	for _, item := range scanResult {
		list.AddItem(
			item.GetSSID(),
			item.GetQuality(),
			1,
			func() {
				cancel()

				form := tview.NewForm().
					AddInputField("SSID", item.GetSSID(), 20, nil, nil).
					AddPasswordField("Password", "", 20, '*', nil).
					AddButton("Connect", nil)
				frameForm := frameDefault(form)
				app.SetRoot(frameForm, true)
			},
		)
	}
}
