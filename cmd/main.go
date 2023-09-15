package main

import (
	"github.com/mr-chelyshkin/rpi4_network_controller/internal/controller"
)

func main() {
	if err := controller.Run(); err != nil {
		panic(err)
	}
}
