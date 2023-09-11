package controller

import (
	"time"

	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/wifi"

	"github.com/rivo/tview"
)

func cmdConnect(app *tview.Application) {
	networks := tview.NewList()
	app.SetRoot(frameDefault(networks), true)

	wifiController, err := wifi.NewWifi()
	if err != nil {
	}
	exec(wifiController, networks)

	go func() {
		ticker := time.NewTicker(rpi4_network_controller.ScanTimeoutSec * time.Second)
		for {
			select {
			case <-ticker.C:
				app.QueueUpdateDraw(
					func() {
						networks.Clear()
						exec(wifiController, networks)
					},
				)
			}
		}
	}()
}

func exec(wifi wifi.Wifi, list *tview.List) {
	scanResult := wifi.Scan()

	for _, item := range scanResult {
		list.AddItem(
			item.GetSSID(),
			item.GetQuality(),
			1,
			nil,
		)
	}
}
