package controller

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func frameDefault(obj tview.Primitive) *tview.Frame {
	return tview.NewFrame(obj).
		AddText("Press Ctrl-C to exit, ESC back to menu", true, tview.AlignLeft, tcell.ColorGray).
		SetBorders(1, 1, 1, 1, 2, 2)
}
