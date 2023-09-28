package controller

import (
	"context"

	"github.com/mr-chelyshkin/rpi4_network_controller/pkg/wifi"
)

// Controller object.
type Controller struct {
	output chan string
}

// New return Controller object.
func New(output chan string) *Controller {
	return &Controller{
		output: output,
	}
}

// Scan scans for available networks and returns the result.
func (c *Controller) Scan(ctx context.Context) []*wifi.Network {
	resultCh := make(chan []*wifi.Network, 1)
	go func() {
		defer close(resultCh)
		resultCh <- wifi.Scan(c.output)
	}()
	select {
	case <-ctx.Done():
		return nil
	case result := <-resultCh:
		return result
	}
}

// Connect tries to connect to a network and returns the result.
func (c *Controller) Connect(ctx context.Context, ssid, password string) bool {
	resultCh := make(chan bool, 1)
	go func() {
		defer close(resultCh)
		resultCh <- wifi.Conn(ssid, password, c.output)
	}()
	select {
	case <-ctx.Done():
		return false
	case result := <-resultCh:
		return result
	}
}

// Status gets the wifi connection status.
func (c *Controller) Status(ctx context.Context) string {
	resultCh := make(chan string, 1)
	go func() {
		defer close(resultCh)
		resultCh <- wifi.State(c.output)
	}()
	select {
	case <-ctx.Done():
		return ""
	case result := <-resultCh:
		return result
	}
}
