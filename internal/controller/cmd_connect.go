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
	go scanner(ctx, cancel, app)
}

func scanner(ctx context.Context, cancel context.CancelFunc, app *tview.Application) {
	controller, _ := wifi.NewWifi()
	scanResults := tview.NewList()

	scan(cancel, controller, app, scanResults)
	scanTick := func() {
		scanResults.Clear()
		scan(cancel, controller, app, scanResults)
	}
	app.SetRoot(frameDefault(scanResults), true)

	ticker := time.NewTicker(rpi4_network_controller.ScanTimeoutSec * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			app.QueueUpdateDraw(scanTick)
		case <-ctx.Done():
			return
		}
	}
}

func scan(cancel context.CancelFunc, controller *wifi.Wifi, app *tview.Application, scanList *tview.List) {
	for _, network := range controller.Scan() {
		network := network

		networkForm := func() {
			cancel()

			writer := tview.NewTextView().
				SetDynamicColors(true).
				SetChangedFunc(func() { app.Draw() })
			form := tview.NewForm().
				AddInputField("SSID", network.GetSSID(), 40, nil, nil).
				AddPasswordField("Password", "", 40, '*', nil)
			form.AddButton(
				"Connect",
				func() {
					conn(
						network,
						controller,
						writer,
						form.GetFormItem(1).(*tview.InputField).GetText(),
					)
				},
			)
			grid := tview.NewGrid().
				SetRows(-5, 1, 0).
				SetBorders(true).
				SetColumns(0)
			grid.AddItem(form, 0, 0, 1, 3, 0, 0, true)
			grid.AddItem(writer, 2, 0, 1, 3, 0, 0, false)

			frameFrom := frameDefault(grid)
			app.SetRoot(frameFrom, true)
		}

		subStr := fmt.Sprintf(
			"Freq: %s | Level: %s | Quality: %s",
			network.GetFreq(),
			network.GetLevel(),
			network.GetQuality(),
		)
		scanList.AddItem(network.GetSSID(), subStr, '*', networkForm)
	}
}

func conn(network *wifi.Network, controller *wifi.Wifi, writer *tview.TextView, password string) {
	logs := make(chan string, 5)
	if len(password) < 8 {
		logs <- "Error: WiFi password should be 8 more chars."
		return
	}

	go func() {
		logs <- fmt.Sprintf("Info: Try connecting to '%s'", network.GetSSID())
		_ = controller.Conn(network.GetSSID(), password, logs)
		logs <- fmt.Sprintf("OK: %s", controller.Active())
	}()
	go func() {
		for {
			select {
			case log := <-logs:
				fmt.Fprintf(writer, "%s\n", log)
			}
		}
	}()
}
