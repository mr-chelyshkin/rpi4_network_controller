package controller

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/rivo/tview"
)

func frameDefault(
	ctx context.Context,
	obj tview.Primitive,
	log chan string,
) *tview.Frame {
	logWriter := tview.NewTextView().
		ScrollToEnd().
		SetDynamicColors(true).
		SetChangedFunc(func() { rpi4_network_controller.App.Draw() })
	grid := tview.NewGrid().
		SetRows(-5, 1, 0).
		SetBorders(true).
		SetColumns(0)
	grid.AddItem(obj, 0, 0, 1, 3, 0, 0, true)
	grid.AddItem(logWriter, 2, 0, 1, 3, 0, 0, false)

	go func() {
		for {
			select {
			case l := <-log:
				fmt.Fprintf(logWriter, "%s\n", l)
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
