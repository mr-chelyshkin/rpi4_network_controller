package main

import (
	"github.com/mr-chelyshkin/rpi4_network_controller/internal"
)

func main() {
	app, err := internal.NewApp()
	if err != nil {
		panic(err)
	}
	app.Run()
}
