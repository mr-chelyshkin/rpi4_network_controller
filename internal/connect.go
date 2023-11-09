package internal

import (
	"context"
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/schedule"
)

func connect(ctx context.Context, interrupt chan struct{}) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	networks := make(chan []map[string]string)

	schedule.NetworkScan(ctx, networks)
	for {
		select {
		case <-networks:
			ctx.Value(rpi4_network_controller.CtxKeyOutputCh).(chan string) <- "networks"
		case <-interrupt:
			return
		case <-ctx.Done():
			return
		}
	}
}
