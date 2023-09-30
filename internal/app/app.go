package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type cmdConnectNetworkDetails struct {
	form     func()
	subTitle string
	title    string
}

// Run application.
func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan struct{}, 1)
	defer close(stop)
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sign:
			cancel()
		case <-stop:
			cancel()
		}
	}()

	frame := tview.NewList().
		AddItem("Connect", "connect to wifi network", '1', func() { cmdConnect(stop) }).
		AddItem("Disconnect", "interrupt wifi connection", '2', func() {})

	app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				setFrame(frameWrapper(ctx, frame, nil))
				stop <- struct{}{}
			case tcell.KeyCtrlC:
				stop <- struct{}{}
				app.Stop()
			}
			return event
		},
	)
	return appRun(frameWrapper(ctx, frame, nil))
}
