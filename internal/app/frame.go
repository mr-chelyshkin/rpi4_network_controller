package app

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()

func frameDraw(frame *tview.Frame) {
	app.SetRoot(frame, true).SetFocus(frame).Draw()
}

func frameWrapper(ctx context.Context, p tview.Primitive, o chan string) *tview.Frame {
	writer := tview.NewTextView().
		ScrollToEnd().
		SetDynamicColors(true).
		SetChangedFunc(func() { app.Draw() })
	grid := tview.NewGrid().
		SetRows(-5, 1, 0).
		SetBorders(true).
		SetColumns(0)
	grid.AddItem(p, 0, 0, 1, 3, 0, 0, true)
	grid.AddItem(writer, 2, 0, 1, 3, 0, 0, false)

	go func() {
		for {
			select {
			case output := <-o:
				_, _ = fmt.Fprintf(writer, "%s\n", output)
			case <-ctx.Done():
				return
			}
		}
	}()
	return tview.NewFrame(grid).
		AddText(
			"Press Ctrl-C to exit, ESC back to menu",
			true,
			tview.AlignLeft,
			tcell.ColorGray,
		).SetBorders(1, 1, 1, 1, 2, 2)
}
