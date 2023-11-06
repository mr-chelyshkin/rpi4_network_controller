package ui

import (
	"context"
	"github.com/rivo/tview"
)

func Frame(ctx context.Context, p tview.Primitive) *tview.Frame {
	content := flexContent(p)
	header := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(flexInfo(ctx), 0, 1, false).
		AddItem(flexHotKeys(ctx), 0, 2, false)
	footer := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(flexWriter(ctx), 0, 1, false)

	frame := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 6, 1, false).
		AddItem(content, 0, 3, true).
		AddItem(footer, 4, 2, false)

	return tview.NewFrame(frame)
}
