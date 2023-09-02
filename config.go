package rpi4_network_controller

import (
	"github.com/gdamore/tcell/v2"
)

const (
	ScanTimeoutSec = 4
)

var (
	AttentionStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorCadetBlue)
	BaseStyle      = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	NoteStyle      = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGrey)
)
