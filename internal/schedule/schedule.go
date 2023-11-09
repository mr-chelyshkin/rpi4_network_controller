package schedule

import (
	"context"
	"os/user"
	"time"

	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/controller"
)

// ADD: ctx withTimeout for jobs.

func UserInfo(ctx context.Context, c chan<- [2]string) {
	var (
		uid = "error"
		usr = "error"
	)

	f := func() {
		u, err := user.Current()
		if err == nil {
			usr = u.Username
			uid = u.Uid
		}
		c <- [2]string{usr, uid}
	}
	go schedule(ctx, rpi4_network_controller.ScanTickGlobal, f)
}

func NetworkStatus(ctx context.Context, c chan<- string) {
	wifi, ok := ctx.Value(rpi4_network_controller.CtxKeyWifiController).(controller.Controller)
	if !ok {
		return
	}

	f := func() { c <- wifi.Status(ctx, nil) }
	go schedule(ctx, rpi4_network_controller.ScanTickGlobal, f)
}

func NetworkScan(ctx context.Context, c chan<- []map[string]string) {
	wifi, ok := ctx.Value(rpi4_network_controller.CtxKeyWifiController).(controller.Controller)
	if !ok {
		return
	}
	output, ok := ctx.Value(rpi4_network_controller.CtxKeyOutputCh).(chan string)
	if !ok {
		return
	}

	f := func() {
		ctx.Value(rpi4_network_controller.CtxKeyOutputCh).(chan string) <- "Run Tick"
		var networks []map[string]string

		for _, network := range wifi.Scan(ctx, output) {
			networks = append(networks, map[string]string{
				"ssid":    network.GetSSID(),
				"quality": network.GetQuality(),
				"freq":    network.GetFreq(),
				"level":   network.GetLevel(),
			})
		}
		c <- networks
	}
	go schedule(ctx, 4, f)
}

func schedule(ctx context.Context, tick int, f func()) {
	ticker := time.NewTicker(time.Duration(tick) * time.Second)
	defer ticker.Stop()

	f()
	for {
		select {
		case <-ticker.C:
			f()
		case <-ctx.Done():
			return
		}
	}
}
