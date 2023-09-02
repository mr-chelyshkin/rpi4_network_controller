package internal

import (
	"github.com/gdamore/tcell/v2"
)

type text struct {
	msg string
}

// InitText ...
func InitText(msg string) *text {
	return &text{msg: msg}
}

// ToScreen ...
func (t text) ToScreen(x, y int, screen tcell.Screen, style tcell.Style) {
	for _, r := range []rune(t.msg) {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}
