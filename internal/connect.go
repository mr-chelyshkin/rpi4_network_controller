package internal

import (
	"context"
	"fmt"
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/schedule"
)

func connect(ctx context.Context, interrupt chan struct{}) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	networks := make(chan []map[string]string, 1)
	schedule.NetworkScan(ctx, networks)
	for {
		select {
		case d := <-networks:
			for _, i := range d {
				ff := fmt.Sprintf("%s, %s", i["ssid"], i["level"])
				ctx.Value(rpi4_network_controller.CtxKeyOutputCh).(chan string) <- ff
			}
			ctx.Value(rpi4_network_controller.CtxKeyOutputCh).(chan string) <- "DONE"
		case <-interrupt:
			return
		case <-ctx.Done():
			return
		}
	}
}
