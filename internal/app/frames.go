package app

import (
	"context"

	"github.com/rivo/tview"
)

var app = tview.NewApplication()

func frameMain(ctx context.Context) *tview.Frame {
	frame := tview.NewList().
		AddItem("Connect", "connect to wifi network", '1', func() { cmdConnect(ctx) }).
		AddItem("Disconnect", "interrupt wifi connection", '2', func() {})

	return frameDefault(ctx, frame, nil)
}

func frameList(ctx context.Context, stdlog chan string) *tview.Frame {
	frame := tview.NewList()

	return frameDefault(ctx, frame, stdlog)
}

func focus(frame *tview.Frame) {
	app.SetRoot(frame, true)
}
