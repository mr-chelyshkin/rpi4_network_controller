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
	go scanner(ctx)

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

func scanner(ctx context.Context) {
	networks := make(chan []cmdConnectNetworkDetails, 1)
	defer close(networks)

	output := make(chan string, 1)
	defer close(output)

	output <- "start scanner: refresh every 4s."

	view := tview.NewList()
	wifi := controller.New()
	frameDraw(frameWrapper(ctx, view, output))
	go func() {
		for {
			select {
			case networks := <-networks:
				app.QueueUpdateDraw(func() {
					view.Clear()

					for _, network := range networks {
						network := network
						view.AddItem(network.title, network.subTitle, '*',
							func() {
								output <- "stop scanner."
								network.form()
							},
						)
					}
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	scan(ctx, output, wifi, networks)

	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			scan(ctx, output, wifi, networks)
		case <-ctx.Done():
			return
		}
	}
}

func scan(
	ctx context.Context,
	output chan string,
	wifi *controller.Controller,
	data chan []cmdConnectNetworkDetails,
) {
	networks := []cmdConnectNetworkDetails{}

	for _, network := range wifi.Scan(ctx, output) {
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
			go conn(
				ctx,
				output,
				network,
				wifi,
				form.GetFormItem(1).(*tview.InputField).GetText(),
			)
		},
	)
	setFrame(frameWrapper(ctx, form, output))
}

func conn(
	ctx context.Context,
	output chan string,
	network *wifi.Network,
	wifi *controller.Controller,
	password string,
) {
	output <- fmt.Sprintf("try connect to %s\n", network.GetSSID())
	_ = wifi.Connect(ctx, output, network.GetSSID(), password)
	output <- wifi.Status(ctx, output)
}
