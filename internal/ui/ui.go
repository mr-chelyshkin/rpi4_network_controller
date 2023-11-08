package ui

import (
	"context"

	"github.com/rivo/tview"
)

var App = tview.NewApplication()

func Start(ctx context.Context, p tview.Primitive) error {
	return setFrame(Main(ctx, p)).Run()
}

func Draw(ctx context.Context, p tview.Primitive) {
	setFrame(Main(ctx, p)).Draw()
}

func setFrame(frame *tview.Frame) *tview.Application {
	return App.SetRoot(frame, true).SetFocus(frame)
}
