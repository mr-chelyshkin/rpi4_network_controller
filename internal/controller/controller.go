package controller

import (
	"context"
	"sort"
	"strings"

	"github.com/mr-chelyshkin/rpi4_network_controller/pkg/wifi"
)

// Controller object.
type Controller struct {
	scanSkipEmptySsid   bool
	scanSortBySignalLvl bool
	scanSortBySsidName  bool
}

// New return Controller object.
func New(opts ...ControllerOpts) Controller {
	c := &Controller{}
	for _, opt := range opts {
		opt(c)
	}
	return *c
}

// Scan available networks and returns the result.
func (c Controller) Scan(ctx context.Context, output chan string) []wifi.Network {
	resultCh := make(chan []wifi.Network, 1)
	go func() {
		defer close(resultCh)

		resultCh <- func() []wifi.Network {
			networks := []wifi.Network{}
			defer func() { networks = nil }()

			for _, network := range wifi.Scan(output) {
				if c.scanSkipEmptySsid && len(network.GetSSID()) == 0 {
					continue
				}
				networks = append(networks, *network)
			}

			switch {
			case c.scanSortBySignalLvl:
				sort.Sort(wifi.ByLevelDesc(networks))
			case c.scanSortBySsidName:
				sort.Sort(wifi.BySsidDesc(networks))
			}
			return networks
		}()
	}()
	select {
	case <-ctx.Done():
		return nil
	case result := <-resultCh:
		return result
	}
}

// Connect tries to connect to a network and returns the result.
func (c Controller) Connect(ctx context.Context, output chan string, ssid, pass, country string) bool {
	resultCh := make(chan bool, 1)
	go func() {
		defer close(resultCh)

		if len(pass) != 0 && len(pass) < 8 {
			output <- "error: WiFi password should be 8 or more chars."
			return
		}
		if len(country) == 0 {
			country = "US"
		}
		resultCh <- wifi.Conn(ssid, pass, strings.ToUpper(country), output)
	}()
	select {
	case <-ctx.Done():
		return false
	case result := <-resultCh:
		return result
	}
}

// Status gets the WiFi connection status.
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
