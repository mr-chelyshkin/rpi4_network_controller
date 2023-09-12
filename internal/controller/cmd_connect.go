package controller

import (
	"context"
	"fmt"
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
		// Handle this error appropriately; for now, just returning
		return
	}
	exec(cancel, wifiController, app, networks)
	go monitorNetworks(ctx, cancel, wifiController, app, networks)
}

func monitorNetworks(
	ctx context.Context,
	cancel context.CancelFunc,
	wifiController wifi.Wifi,
	app *tview.Application,
	networks *tview.List,
) {
	ticker := time.NewTicker(rpi4_network_controller.ScanTimeoutSec * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			app.QueueUpdateDraw(func() {
				networks.Clear()
				exec(cancel, wifiController, app, networks)
			})
		case <-ctx.Done():
			return
		}
	}
}

func exec(cancel context.CancelFunc, wifi wifi.Wifi, app *tview.Application, list *tview.List) {
	scanResult := wifi.Scan()

	for _, item := range scanResult {
		item := item
		list.AddItem(item.GetSSID(), item.GetQuality(), '*', func() {
			cancel()

			stdoutCH := make(chan string)
			// defer close(stdoutCH)

			form := tview.NewForm().
				AddInputField("SSID", item.GetSSID(), 20, nil, nil).
				AddPasswordField("Password", "", 20, '*', nil)

			form.AddButton("Connect", func() {
				ssid := item.GetSSID()
				pass := form.GetFormItem(1).(*tview.InputField).GetText()
				if len(pass) < 8 {
					stdoutCH <- "passwod to short\n"
					return
				}

				stdoutCH <- fmt.Sprintf("connecting to %s", ssid)
				_ = wifi.Conn(ssid, pass, stdoutCH)
				stdoutCH <- fmt.Sprintf("network status: %s", wifi.Active())
				// close(stdoutCH)
			})

			//
			stdout := tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() { app.Draw() })
			go func() {
				for {
					select {
					case out := <-stdoutCH:
						fmt.Fprintf(stdout, "%s", out)
					}
				}
			}()

			//
			grid := tview.NewGrid().SetRows(-5, 1, 0).SetColumns(0).SetBorders(true)
			grid.AddItem(form, 0, 0, 1, 3, 0, 0, true)
			grid.AddItem(stdout, 2, 0, 1, 3, 0, 0, false)

			frameForm := frameDefault(grid)
			app.SetRoot(frameForm, true)
		})
	}
}
