package schedule

import (
	"context"
	"os/user"
	"time"

	"github.com/mr-chelyshkin/rpi4_network_controller"
)

func UserInfo(ctx context.Context, c chan<- [2]string) {
	f := func() {
		var (
			uid = "error"
			usr = "error"
		)
		u, err := user.Current()
		if err == nil {
			usr = u.Username
			uid = u.Uid
		}
		c <- [2]string{usr, uid}
	}
	schedule(ctx, rpi4_network_controller.ScanTickGlobal, f)
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
