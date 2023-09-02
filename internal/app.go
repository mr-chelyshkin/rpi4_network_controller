package internal

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/pkg/wifi"
)

// App ...
type App struct {
	terminalScreen tcell.Screen
	wifiController wifi.Wifi
	wifiScanResult table
}

// NewApp ...
func NewApp() (*App, error) {
	terminalScreen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := terminalScreen.Init(); err != nil {
		return nil, err
	}
	wifiController, err := wifi.NewWifi()
	if err != nil {
		return nil, err
	}

	return &App{
		terminalScreen: terminalScreen,
		wifiController: wifiController,
	}, nil
}

// Run ...
func (a *App) Run() error {
	ticker := time.NewTicker(rpi4_network_controller.ScanTimeoutSec * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			InitText("'Esc' to exit or 'Enter' for input network number.").ToScreen(1, 0, a.terminalScreen, rpi4_network_controller.NoteStyle)
			InitText(a.wifiController.Active()).ToScreen(1, 2, a.terminalScreen, rpi4_network_controller.AttentionStyle)

			wifiScanResult := InitTable([]string{"#", "SSID", "Level"})
			for idx, network := range a.wifiController.Scan() {
				wifiScanResult.AddRow(
					[]string{
						strconv.Itoa(idx),
						network.GetSSID(),
						network.GetLevel(),
					})
			}
			wifiScanResult.ToScreen(1, 4, a.terminalScreen, rpi4_network_controller.BaseStyle)
			a.wifiScanResult = wifiScanResult

			a.terminalScreen.Sync()
			fmt.Print("\x1b[3J")
		}
	}()
	for {
		ev := a.terminalScreen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			a.terminalScreen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				a.terminalScreen.Fini()
				return nil
			}
			if ev.Key() == tcell.KeyCtrlC {
				a.terminalScreen.Fini()
				return nil
			}
			if ev.Key() == tcell.KeyEnter {
				a.terminalScreen.Fini()

				var index int
				var password string

				fmt.Print("Select network number: ")
				fmt.Scan(&index)

				if index >= 0 && index < len(a.wifiScanResult.rows) {
					fmt.Print("Enter password: ")
					fmt.Scan(&password)

					selectedNetwork := a.wifiScanResult.rows[index]
					fmt.Println(selectedNetwork[1])
					dd := a.wifiController.Conn(selectedNetwork[1], password)
					fmt.Println(dd)
				}

				b := a.wifiController.Active()
				fmt.Println(b)
			}
		}
	}
}
