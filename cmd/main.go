package main

import (
	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/app"
	"os/user"
)

func init() {
	u, err := user.Current()
	if err == nil {
		rpi4_network_controller.UserName = u.Username
		rpi4_network_controller.UserPerm = u.Uid
	}
}

func main() {
	app.Run()
}
