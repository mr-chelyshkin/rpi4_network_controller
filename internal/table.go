package internal

import (
	"github.com/gdamore/tcell/v2"
)

type table struct {
	headers []string
	rows    [][]string
}

func InitTable(headers []string) table {
	return table{headers: headers}
}

// AddRow ...
func (t *table) AddRow(row []string) {
	t.rows = append(t.rows, row)
}

// ToScreen ...
func (t *table) ToScreen(x, y int, screen tcell.Screen, style tcell.Style) {
	colWidths := make([]int, len(t.headers))
	for col, header := range t.headers {
		colWidths[col] = len(header)
		for _, row := range t.rows {
			if len(row[col]) > colWidths[col] {
				colWidths[col] = len(row[col])
			}
		}
	}
	for col, header := range t.headers {
		for _, ch := range header {
			screen.SetContent(x, y, ch, nil, style.Underline(true))
			x++
		}
		x += colWidths[col] - len(header) + 1
	}

	y += 2
	for _, row := range t.rows {
		x = 1
		for col, cell := range row {
			for _, ch := range cell {
				screen.SetContent(x, y, ch, nil, style)
				x++
			}
			x += colWidths[col] - len(cell) + 1
		}
		y += 2
	}
}
