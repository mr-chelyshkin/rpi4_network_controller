package internal

import (
	"context"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/ui"
)

// Mutex on commend

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan struct{}, 1)
	defer close(stop)

	output := make(chan string, 1)
	ctx = context.WithValue(ctx, rpi4_network_controller.CtxKeyOutputCh, output)
	ctx = context.WithValue(ctx, rpi4_network_controller.CtxKeyHotkeys, []ui.HotKeys{
		{
			Key:         tcell.KeyESC,
			Description: "Go to main menu",
			Action: func(ctx context.Context) {
				output <- "send interrupt"
				stop <- struct{}{}
				//Run()
			},
		},
		{
			Key:         tcell.KeyCtrlC,
			Description: "Exit",
			Action: func(ctx context.Context) {
				stop <- struct{}{}
				ui.App.Stop()
				os.Exit(0)
			},
		},
	})
	ctx = context.WithValue(ctx, rpi4_network_controller.CtxKeyWifiController, controller.New(
		controller.WithScanSkipEmptySSIDs(),
		controller.WithScanSortByLevel(),
	))
	view := ui.ContentTable(ctx, ui.ContentTableData{
		Headers: []string{"action", "description"},
		Data: []ui.ContentTableRow{
			{
				Action: func(ctx context.Context) { go connect(ctx, stop) },
				Data:   []string{"connect", "scan and connect to wifi network"},
			},
			{
				Action: func(ctx context.Context) {},
				Data:   []string{"disconnect", "interrupt current wifi connection"},
			},
		},
	})
	return ui.Start(ctx, view)
}
