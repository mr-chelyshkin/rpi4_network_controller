package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/ui"

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
	ctx = context.WithValue(ctx, rpi4_network_controller.CtxKeyOutputCh, output)
	defer close(output)

	output <- "start scanner: refresh every 4s."

	view := tview.NewList()
	wifi := controller.New(
		controller.WithScanSkipEmptySSIDs(),
		controller.WithScanSortByLevel(),
	)
	//ui.Draw(frameWrapper(ctx, view, output))
	ui.Draw(ctx, view)
	go func() {
		for {
			select {
			case networks := <-networks:
				ui.App.QueueUpdateDraw(func() {
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
	scan(ctx, cancel, output, wifi, networks)

	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			scan(ctx, cancel, output, wifi, networks)
		case <-ctx.Done():
			return
		}
	}
}

func scan(
	ctx context.Context,
	cancel context.CancelFunc,
	output chan string,
	wifi controller.Controller,
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
				form: func() {
					cancel()
					connForm(network, wifi)
				},
			},
		)
	}
	data <- networks
}

func connForm(network wifi.Network, wifi controller.Controller) {
	ctx := context.Background()
	output := make(chan string, 1)

	form := tview.NewForm().
		AddInputField("SSID", network.GetSSID(), 40, nil, nil).
		AddInputField("Country code", "US", 40, nil, nil).
		AddPasswordField("Password", "", 40, '*', nil)
	form.AddButton(
		"Connect",
		func() {
			go conn(
				ctx,
				output,
				wifi,
				form.GetFormItem(0).(*tview.InputField).GetText(),
				form.GetFormItem(1).(*tview.InputField).GetText(),
				form.GetFormItem(2).(*tview.InputField).GetText(),
			)
		},
	)
	ui.Draw(ctx, form)
	//ui.Draw(frameWrapper(ctx, form, output))
}

func conn(
	ctx context.Context,
	output chan string,
	wifi controller.Controller,
	ssid,
	country,
	pass string,
) {
	output <- fmt.Sprintf("try connect to %s\n", ssid)
	_ = wifi.Connect(ctx, output, ssid, pass, country)
	output <- wifi.Status(ctx, output)
}
