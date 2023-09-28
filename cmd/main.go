package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/controller"

	"github.com/rivo/tview"
)

func cmdConnect(ctx context.Context, stop chan struct{}) {
	cctx, cancel := context.WithCancel(ctx)
	go scanner(cctx, cancel)

	go func() {
		for {
			select {
			case <-stop:
				cancel()
				return
			}
		}
	}()
}

func scanner(ctx context.Context, cancel context.CancelFunc) {
	scanResults := tview.NewList()
	log := make(chan string, 1)

	controller := controller.New(log)
	rpi4_network_controller.App.SetRoot(
		frameDefault(ctx, scanResults, log),
		true,
	)
	log <- fmt.Sprintf(
		"Update networks every %d(sec)\n",
		rpi4_network_controller.ScanTimeoutSec,
	)

	ticker := time.NewTicker(
		rpi4_network_controller.ScanTimeoutSec * time.Second,
	)
	defer ticker.Stop()
	scan(ctx, cancel, controller, scanResults)
	for {
		select {
		case <-ticker.C:
			rpi4_network_controller.App.QueueUpdateDraw(
				func() {
					scan(ctx, cancel, controller, scanResults)
				},
			)
		case <-ctx.Done():
			close(log)
			return
		}
	}
}

func scan(
	ctx context.Context,
	cancel context.CancelFunc,
	controller *wifi1.Wifi,
	scanList *tview.List,
) {
	type networkDetails struct {
		description string
		form        func()
		title       string
	}
	networks := []networkDetails{}
	refresh := make(chan struct{}, 1)

	go func() {
		for _, network := range controller.Scan() {
			network := network

			networkForm := func() {
				cancel()
				form(network, controller)
			}

			description := fmt.Sprintf(
				"Freq: %s | Level: %s | Quality: %s",
				network.GetFreq(),
				network.GetLevel(),
				network.GetQuality(),
			)
			networks = append(
				networks,
				networkDetails{
					form:        networkForm,
					description: description,
					title:       network.GetSSID(),
				},
			)
		}
		refresh <- struct{}{}
	}()
	go func() {
		for {
			select {
			case <-refresh:
				scanList.Clear()

				for _, network := range networks {
					scanList.AddItem(
						network.title,
						network.description,
						'*',
						network.form,
					)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func form(
	network *wifi1.Network,
	controller *wifi1.Wifi,
) {
	ctx := context.Background()
	log := make(chan string, 1)

	form := tview.NewForm().
		AddInputField("SSID", network.GetSSID(), 40, nil, nil).
		AddPasswordField("Password", "", 40, '*', nil)
	form.AddButton(
		"Connect",
		func() {
			conn(
				log,
				network,
				controller,
				form.GetFormItem(1).(*tview.InputField).GetText(),
			)
		},
	)
	frameForm := frameDefault(ctx, form, log)
	rpi4_network_controller.App.SetRoot(frameForm, true)
}

func conn(
	log chan string,
	network *wifi1.Network,
	controller *wifi1.Wifi,
	password string,
) {
	log <- fmt.Sprintf("Try connect to %s", network.GetSSID())
	if len(password) < 8 {
		log <- "WiFi password should be 8 or more chars"
		return
	}

	go func() {
		_ = controller.Conn(network.GetSSID(), password, log)
		log <- controller.Active()
	}()
}
