package main

import (
	"os/user"

	"github.com/mr-chelyshkin/rpi4_network_controller"
	"github.com/mr-chelyshkin/rpi4_network_controller/internal"
)

func init() {
	u, err := user.Current()
	if err == nil {
		rpi4_network_controller.UserName = u.Username
		rpi4_network_controller.UserPerm = u.Uid
	}
}

func main() {
	internal.Run()
}
