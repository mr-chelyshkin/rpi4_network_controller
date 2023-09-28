package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mr-chelyshkin/rpi4_network_controller/internal/controller"

	"github.com/rivo/tview"
)

type networkDetails struct {
	description string
	form        func()
	title       string
}

func cmdConnect(interrupt chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	go scanner(ctx, cancel)

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-interrupt:
			cancel()
			return
		}
	}()
}

func scanner(ctx context.Context, cancel context.CancelFunc) {
	networks := make(chan []networkDetails, 1)
	defer close(networks)

	output := make(chan string, 1)
	defer close(output)

	view := tview.NewList()
	wifi := controller.New(output)
	frameDraw(frameWrapper(ctx, view, output))
	go func() {
		for {
			select {
			case networks := <-networks:
				app.QueueUpdateDraw(func() {
					view.Clear()

					for _, network := range networks {
						view.AddItem(
							network.title,
							network.description,
							'*',
							nil,
						)
					}
				})
			case <-ctx.Done():
				return
			}
		}
	}()

	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			scan(ctx, networks, wifi)
		case <-ctx.Done():
			return
		}
	}
}

func scan(ctx context.Context, data chan []networkDetails, wifi *controller.Controller) {
	networks := []networkDetails{}

	for _, network := range wifi.Scan(ctx) {
		network := network

		description := fmt.Sprintf(
			"Freq: %s | Level: %s | Quality: %s",
			network.GetFreq(),
			network.GetLevel(),
			network.GetQuality(),
		)
		networks = append(
			networks,

			networkDetails{
				description: description,
				title:       network.GetSSID(),
			},
		)
	}
	data <- networks
}

//func form(
//	network *wifi1.Network,
//	controller *wifi1.Wifi,
//) {
//	ctx := context.Background()
//	log := make(chan string, 1)
//
//	form := tview.NewForm().
//		AddInputField("SSID", network.GetSSID(), 40, nil, nil).
//		AddPasswordField("Password", "", 40, '*', nil)
//	form.AddButton(
//		"Connect",
//		func() {
//			conn(
//				log,
//				network,
//				controller,
//				form.GetFormItem(1).(*tview.InputField).GetText(),
//			)
//		},
//	)
//	frameForm := frameDefault(ctx, form, log)
//	rpi4_network_controller.App.SetRoot(frameForm, true)
//}

//func conn(
//	log chan string,
//	network *wifi1.Network,
//	controller *wifi1.Wifi,
//	password string,
//) {
//	log <- fmt.Sprintf("Try connect to %s", network.GetSSID())
//	if len(password) < 8 {
//		log <- "WiFi password should be 8 or more chars"
//		return
//	}
//
//	go func() {
//		_ = controller.Conn(network.GetSSID(), password, log)
//		log <- controller.Active()
//	}()
//}
