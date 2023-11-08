package ui

import (
	"context"

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
