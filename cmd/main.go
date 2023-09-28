package main

import (
	"fmt"

	"github.com/mr-chelyshkin/rpi4_network_controller/pkg/wifi"
)

func main() {
	ch := make(chan string, 1)

	res := wifi.State(ch)
	go func() {
		for {
			select {
			case l := <-ch:
				fmt.Println(l)
			}
		}
	}()
	fmt.Println("end")
	fmt.Println(res)

	//if err := controller.Run(); err != nil {
	//	panic(err)
	//}
}
