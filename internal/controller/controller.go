package controller

import (
	"context"

	"github.com/mr-chelyshkin/rpi4_network_controller/pkg/wifi"
)

type controller struct {
	output chan string
}

func New(output chan string) *controller {
	return &controller{
		output: output,
	}
}

// Scan scans for available networks and returns the result.
func (c *controller) Scan(ctx context.Context) []*wifi.Network {
	resultCh := make(chan []*wifi.Network, 1)
	go func() {
		defer close(resultCh)
		resultCh <- wifi.Scan(c.output)
	}()
	select {
	case <-ctx.Done():
		return nil // return nil or handle this situation as you see fit
	case result := <-resultCh:
		return result
	}
}

// Connect tries to connect to a network and returns the result.
func (c *controller) Connect(ctx context.Context, ssid, password string) bool {
	resultCh := make(chan bool, 1)
	go func() {
		defer close(resultCh)
		resultCh <- wifi.Conn(ssid, password, c.output)
	}()
	select {
	case <-ctx.Done():
		return false // return false or handle this situation as you see fit
	case result := <-resultCh:
		return result
	}
}

// Status gets the wifi connection status.
func (c *controller) Status(ctx context.Context) string {
	resultCh := make(chan string, 1)
	go func() {
		defer close(resultCh)
		resultCh <- wifi.State(c.output)
	}()
	select {
	case <-ctx.Done():
		return "" // return an empty string or handle this situation as you see fit
	case result := <-resultCh:
		return result
	}
}
