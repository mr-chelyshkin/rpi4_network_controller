package internal

import (
	"context"
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/schedule"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/ui"
)

func connect(ctx context.Context, interrupt chan struct{}) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//data := ui.ContentTableData{Headers: []string{"ssid", "freq", "quality", "level"}}

	networks := make(chan []map[string]string, 1)
	schedule.NetworkScan(ctx, networks)
	for {
		select {
		case networks := <-networks:
			rows := []ui.ContentTableRow{}

			for _, network := range networks {
				rows = append(rows, ui.ContentTableRow{
					Data:   []string{network["ssid"], network["freq"], network["quality"], network["level"]},
					Action: func(ctx context.Context) {},
				})
				ctx.Value(rpi4_network_controller.CtxKeyOutputCh).(chan string) <- network["ssid"]
			}
			////ui.App.QueueUpdateDraw(func() {
			//data.Data = rows
			//ui.Draw(ctx, ui.ContentTable(ctx, data))
			////})
			ctx.Value(rpi4_network_controller.CtxKeyOutputCh).(chan string) <- "Updated"
		case <-interrupt:
			return
		case <-ctx.Done():
			return
		}
	}
}
