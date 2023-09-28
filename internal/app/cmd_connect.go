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
				view.Clear()

				for _, network := range networks {
					view.AddItem(
						network.title,
						network.description,
						'*',
						nil,
					)
				}
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
			//app.QueueUpdateDraw(func() {
			//	view.AddItem("ADDDDD", "AAAAAAAAAAAAA", '*', func() {})
			//})
		case <-ctx.Done():
			return
		}
	}

	//	networksCh := make(chan []networkDetails, 1)
	// defer close(networksCh)

	//	go func() {
	//		scan(ctx, networksCh, c)
	//	}()
	//	select {
	//	case networks := <-networksCh:
	//		for _, network := range networks {
	//			res.AddItem(
	//				network.title,
	//				network.description,
	//				'*',
	//				network.form,
	//			)
	//		}
	//	case <-ctx.Done():
	// close(output)
	//		res.Clear()
	//		return
	//	}
	//}

	// app.QueueUpdateDraw(tick)
	// ticker := time.NewTicker(4 * time.Second)
	// defer ticker.Stop()

	//for {
	//	select {
	//	case <-ticker.C:
	//		fmt.Print("SSSSFSF")
	//		app.QueueUpdateDraw(func() {})
	//	case <-ctx.Done():
	//		close(output)
	//		return
	//	}
	//}
}

func scan(ctx context.Context, n chan []networkDetails, sc *controller.Controller) {
	networks := []networkDetails{}

	for _, network := range sc.Scan(ctx) {
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
	fmt.Print("AAAAAAAAAAAAA")
	n <- networks
}

//
//func scan(ctx context.Context, sc func(ctx context.Context) []*wifi.Network, fr *tview.List) {
//
//	networks := []networkDetails{}
//	refresh := make(chan struct{}, 1)
//
//	go func() {
//		for _, network := range sc(ctx) {
//			network := network
//
//			description := fmt.Sprintf(
//				"Freq: %s | Level: %s | Quality: %s",
//				network.GetFreq(),
//				network.GetLevel(),
//				network.GetQuality(),
//			)
//			networks = append(
//				networks,
//				networkDetails{
//					description: description,
//					title:       network.GetSSID(),
//				},
//			)
//		}
//		refresh <- struct{}{}
//	}()
//	go func() {
//		for {
//			select {
//			case <-refresh:
//				fr.Clear()
//
//				for _, network := range networks {
//					fr.AddItem(
//						network.title,
//						network.description,
//						'*',
//						network.form,
//					)
//				}
//			case <-ctx.Done():
//				return
//			}
//		}
//	}()
//}

//func scan(
//	ctx context.Context,
//	cancel context.CancelFunc,
//	controller *wifi1.Wifi,
//	scanList *tview.List,
//) {
//	refresh := make(chan struct{}, 1)
//
//	go func() {
//		for _, network := range controller.Scan() {
//			network := network
//
//			networkForm := func() {
//				cancel()
//				form(network, controller)
//			}
//
//			description := fmt.Sprintf(
//				"Freq: %s | Level: %s | Quality: %s",
//				network.GetFreq(),
//				network.GetLevel(),
//				network.GetQuality(),
//			)
//			networks = append(
//				networks,
//				networkDetails{
//					form:        networkForm,
//					description: description,
//					title:       network.GetSSID(),
//				},
//			)
//		}
//		refresh <- struct{}{}
//	}()
//	go func() {
//		for {
//			select {
//			case <-refresh:
//				scanList.Clear()
//
//				for _, network := range networks {
//					scanList.AddItem(
//						network.title,
//						network.description,
//						'*',
//						network.form,
//					)
//				}
//			case <-ctx.Done():
//				return
//			}
//		}
//	}()
//}

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
