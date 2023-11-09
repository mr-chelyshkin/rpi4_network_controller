package ui

import (
	"context"

	"github.com/rivo/tview"
)

func Main(ctx context.Context, p tview.Primitive) *tview.Frame {
	content := primitive(p)
	header := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(info(ctx), 0, 1, false).
		AddItem(hotKeys(ctx), 0, 2, false)
	footer := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(writer(ctx), 0, 1, false)
	frame := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 6, 1, false).
		AddItem(content, 0, 3, true).
		AddItem(footer, 16, 4, false)
	return tview.NewFrame(frame)
}
