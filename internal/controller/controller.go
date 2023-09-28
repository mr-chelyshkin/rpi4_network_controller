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

// Scan ...
func (c *controller) Scan(ctx context.Context) []*wifi.Network {
	resultCh := make(chan []*wifi.Network, 1)

	go func() {
		defer close(resultCh)

		go func() {
			resultCh <- wifi.Scan(c.output)
		}()
		select {
		case <-ctx.Done():
			return
		case result := <-resultCh:
			resultCh <- result
		}
	}()
	return <-resultCh
}

// Connect ...
func (c *controller) Connect(ctx context.Context, ssid, password string) bool {
	resultCh := make(chan bool, 1)

	go func() {
		defer close(resultCh)

		go func() {
			resultCh <- wifi.Conn(ssid, password, c.output)
		}()
		select {
		case <-ctx.Done():
			return
		case result := <-resultCh:
			resultCh <- result
		}
	}()
	return <-resultCh
}

// Status ...
func (c *controller) Status(ctx context.Context) string {
	resultCh := make(chan string, 1)

	go func() {
		defer close(resultCh)

		go func() {
			resultCh <- wifi.State(c.output)
		}()
		select {
		case <-ctx.Done():
			return
		case result := <-resultCh:
			resultCh <- result
		}
	}()
	return <-resultCh
}
