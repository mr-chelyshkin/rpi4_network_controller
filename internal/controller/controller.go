package controller

import (
	"context"

	"github.com/mr-chelyshkin/rpi4_network_controller/pkg/wifi"
)

// Controller object.
type Controller struct{}

// New return Controller object.
func New() Controller {
	return Controller{}
}

// Scan scans for available networks and returns the result.
func (c Controller) Scan(ctx context.Context, output chan string) []*wifi.Network {
	resultCh := make(chan []*wifi.Network, 1)
	go func() {
		defer close(resultCh)
		resultCh <- wifi.Scan(output)
	}()
	select {
	case <-ctx.Done():
		return nil
	case result := <-resultCh:
		return result
	}
}

// Connect tries to connect to a network and returns the result.
func (c Controller) Connect(ctx context.Context, output chan string, ssid, password string) bool {
	resultCh := make(chan bool, 1)
	go func() {
		defer close(resultCh)

		if len(password) != 0 && len(password) < 8 {
			output <- "error: WiFi password should be 8 or more chars."
			return
		}
		resultCh <- wifi.Conn(ssid, password, output)
	}()
	select {
	case <-ctx.Done():
		return false
	case result := <-resultCh:
		return result
	}
}

// Status gets the wifi connection status.
func (c Controller) Status(ctx context.Context, output chan string) string {
	resultCh := make(chan string, 1)
	go func() {
		defer close(resultCh)
		resultCh <- wifi.State(output)
	}()
	select {
	case <-ctx.Done():
		return ""
	case result := <-resultCh:
		return result
	}
}
