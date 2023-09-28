package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mr-chelyshkin/rpi4_network_controller/internal/controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/pkg/wifi"

	"github.com/rivo/tview"
)

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
	networks := make(chan []cmdConnectNetworkDetails, 1)
	defer close(networks)

	output := make(chan string, 1)
	defer close(output)

	output <- "start scanner: refresh every 4s."

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
							network.subTitle,
							'*',
							func() {
								output <- "stop scanner."
								cancel()
							},
						)
					}
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	scan(ctx, wifi, networks)

	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			scan(ctx, wifi, networks)
		case <-ctx.Done():
			return
		}
	}
}

func scan(
	ctx context.Context,
	wifi *controller.Controller,
	data chan []cmdConnectNetworkDetails,
) {
	networks := []cmdConnectNetworkDetails{}

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

			cmdConnectNetworkDetails{
				subTitle: description,
				title:    network.GetSSID(),
				form:     func() { connForm(network, wifi) },
			},
		)
	}
	data <- networks
}

func connForm(network *wifi.Network, wifi *controller.Controller) {
	ctx := context.Background()
	output := make(chan string, 1)

	form := tview.NewForm().
		AddInputField("SSID", network.GetSSID(), 40, nil, nil).
		AddPasswordField("Password", "", 40, '*', nil)
	form.AddButton(
		"Connect",
		func() {
			conn(
				ctx,
				output,
				network,
				wifi,
				form.GetFormItem(1).(*tview.InputField).GetText(),
			)
		},
	)
	frameDraw(frameWrapper(ctx, form, output))
}

func conn(
	ctx context.Context,
	output chan string,
	network *wifi.Network,
	wifi *controller.Controller,
	password string,
) {
	output <- fmt.Sprintf("try connect to %s\n", network.GetSSID())
	if len(password) < 8 {
		output <- "error: WiFi password should be 8 or more chars."
		return
	}

	_ = wifi.Connect(ctx, network.GetSSID(), password)
	output <- wifi.Status(ctx)
}
