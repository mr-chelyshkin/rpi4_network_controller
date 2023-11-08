package ui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/schedule"
	"github.com/rivo/tview"
)

func primitive(p tview.Primitive) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p, 0, 1, true)

	flex.SetBorder(true)
	return flex
}

func writer(ctx context.Context) *tview.Flex {
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

type HotKeys struct {
	Action      func(ctx context.Context)
	Description string
	Key         tcell.Key
}

func hotKeys(ctx context.Context) *tview.Flex {
	frame := tview.NewTable()
	content, ok := ctx.Value(rpi4_network_controller.CtxKeyHotkeys).([]HotKeys)
	if ok {
		row := 0
		for _, key := range content {
			frame.SetCell(row, 0, tview.NewTableCell("<"+tcell.KeyNames[key.Key]+">").SetTextColor(tcell.ColorBlue))
			frame.SetCell(row, 1, tview.NewTableCell(key.Description).SetTextColor(tcell.ColorGray))
			row++
		}
	}
	App.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			for _, k := range content {
				if k.Key == event.Key() {
					k.Action(ctx)
					break
				}
			}
			return event
		},
	)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(frame, 0, 1, false)

	flex.SetBorder(false)
	return flex
}

func info(ctx context.Context) *tview.Flex {
	frame := tview.NewTable()
	frame.SetCell(0, 0, tview.NewTableCell("Version:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(0, 1, tview.NewTableCell(rpi4_network_controller.Version).SetTextColor(tcell.ColorWhite))
	frame.SetCell(1, 0, tview.NewTableCell("User:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(1, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorOrangeRed))
	frame.SetCell(2, 0, tview.NewTableCell("Privileged:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(2, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorOrangeRed))
	frame.SetCell(3, 0, tview.NewTableCell("CurrentConn:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(3, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorOrangeRed))

	usrInfoCh := make(chan [2]string, 1)
	go func() {
		go schedule.UserInfo(ctx, usrInfoCh)
		for {
			select {
			case info := <-usrInfoCh:
				App.QueueUpdateDraw(func() {
					switch info[0] {
					case "error":
						frame.SetCell(1, 1, tview.NewTableCell("error").SetTextColor(tcell.ColorRed))
					default:
						frame.SetCell(1, 1, tview.NewTableCell(info[0]).SetTextColor(tcell.ColorWhite))
					}

					switch info[1] {
					case "error":
						frame.SetCell(2, 1, tview.NewTableCell("error").SetTextColor(tcell.ColorRed))
					case "0":
						frame.SetCell(2, 1, tview.NewTableCell("yes").SetTextColor(tcell.ColorWhite))
					default:
						frame.SetCell(2, 1, tview.NewTableCell("run app with privileged mode").SetTextColor(tcell.ColorRed))
					}
				})
			case <-ctx.Done():
				close(usrInfoCh)
				return
			}
		}
	}()
	networkStatusCh := make(chan string, 1)
	go func() {
		go schedule.NetworkStatus(ctx, networkStatusCh)
		for {
			select {
			case info := <-networkStatusCh:
				App.QueueUpdateDraw(func() {
					frame.SetCell(3, 1, tview.NewTableCell(info).SetTextColor(tcell.ColorOrangeRed))
				})
			case <-ctx.Done():
				close(networkStatusCh)
				return
			}
		}
	}()

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(frame, 0, 1, false)
	flex.SetBorder(false)
	return flex
}
