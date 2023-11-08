package ui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/rivo/tview"
)

type ContentTableRow struct {
	Action func(ctx context.Context)
	Data   []string
}

type ContentTableData struct {
	Headers []string
	Data    []ContentTableRow
}

func ContentTable(ctx context.Context, data ContentTableData) *tview.Table {
	content := tview.NewTable().SetSelectable(true, false)
	for i, header := range data.Headers {
		content.SetCell(0, i, tview.NewTableCell(header).SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	for r, row := range data.Data {
		for c, col := range row.Data {
			content.SetCell(r+1, c, tview.NewTableCell(col))
		}
	}
	content.SetSelectedFunc(func(r, c int) {
		data.Data[r-1].Action(ctx)
	})
	return content
}

func flexHotKeys(ctx context.Context) *tview.Flex {
	frame := tview.NewTable()
	content, ok := ctx.Value(rpi4_network_controller.CtxKeyHotkeys).(map[string]string)
	if ok {
		row := 0
		for k, v := range content {
			frame.SetCell(row, 0, tview.NewTableCell("<"+k+">").SetTextColor(tcell.ColorBlue))
			frame.SetCell(row, 1, tview.NewTableCell(v).SetTextColor(tcell.ColorGray))
			row++
		}
	}
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(frame, 0, 1, false)

	flex.SetBorder(false)
	return flex
}

func flexInfo(ctx context.Context) *tview.Flex {
	frame := tview.NewTable()
	frame.SetCell(0, 0, tview.NewTableCell("Version:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(0, 1, tview.NewTableCell(rpi4_network_controller.Version).SetTextColor(tcell.ColorWhite))
	frame.SetCell(1, 0, tview.NewTableCell("User:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(1, 1, tview.NewTableCell(rpi4_network_controller.UserName).SetTextColor(tcell.ColorWhite))
	frame.SetCell(2, 0, tview.NewTableCell("Privileged:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(2, 1, func() *tview.TableCell {
		if rpi4_network_controller.UserPerm == "0" {
			return tview.NewTableCell("yes").SetTextColor(tcell.ColorWhite)
		}
		return tview.NewTableCell("run app with privileged mode").SetTextColor(tcell.ColorRed)
	}())
	frame.SetCell(3, 0, tview.NewTableCell("CurrentConn:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(3, 1, func() *tview.TableCell {
		content, ok := ctx.Value(rpi4_network_controller.CtxKeyCurConn).(string)
		if ok {
			return tview.NewTableCell(content).SetTextColor(tcell.ColorWhite)
		}
		return tview.NewTableCell("n/a").SetTextColor(tcell.ColorOrange)
	}())

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(frame, 0, 1, false)

	flex.SetBorder(false)
	return flex
}

func flexContent(p tview.Primitive) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p, 0, 1, true)

	flex.SetBorder(true)
	return flex
}

func flexWriter(ctx context.Context) *tview.Flex {
	content, ok := ctx.Value(rpi4_network_controller.CtxKeyOutputCh).(chan string)
	if !ok {
		return tview.NewFlex()
	}
	frame := tview.NewTextView().
		SetChangedFunc(func() { App.Draw() }).
		SetDynamicColors(true).
		ScrollToEnd()
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(frame, 0, 1, false)

	go func() {
		for {
			select {
			case output := <-content:
				App.QueueUpdateDraw(func() {
					_, _ = fmt.Fprintf(frame, "%s\n", output)
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	flex.SetBorder(false)
	return flex
}
