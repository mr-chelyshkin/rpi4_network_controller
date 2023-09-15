//go:build linux
// +build linux

package wifi

/*
#cgo linux CFLAGS: -Dlinux
#cgo linux LDFLAGS: -liw
#include "linux_wifi.h"
*/
import "C"

import (
	"fmt"
	"sort"
	"unsafe"
)

var outputChan chan string

//export goSendToChannel
func goSendToChannel(s *C.char) {
	outputChan <- C.GoString(s)
}

// Wifi ...
type Wifi struct{}

// NewWifi ...
func NewWifi() (*Wifi, error) {
	return &Wifi{}, nil
}

// Scan ...
func (w *Wifi) Scan() []*Network {
	uniqueMap := make(map[string]struct{})
	var uniqueRes []*Network

	for _, net := range scanCGO() {
		if C.GoString(&net.sSID[0]) == "" {
			continue
		}

		key := C.GoString(&net.sSID[0]) + "_" + fmt.Sprintf("%f", net.freq)
		if _, found := uniqueMap[key]; !found {
			uniqueRes = append(uniqueRes, net)
			uniqueMap[key] = struct{}{}
		}
	}
	sort.Sort(byLevelDesc(uniqueRes))
	return uniqueRes
}

// Conn ...
func (w *Wifi) Conn(ssid, password string, output chan string) bool {
	return connCGO(ssid, password, output)
}

// Active ...
func (w *Wifi) Active() string {
	res := activeCGO()

	switch res {
	case "":
		return "No connection"
	default:
		return fmt.Sprintf("Current network: %s", res)
	}
}

func scanCGO() []*Network {
	count := C.int(0)
	results := C.network_scan(&count)
	networks := make([]*Network, count)

	for i := 0; i < int(count); i++ {
		n := (*Network)(
			unsafe.Pointer(
				uintptr(unsafe.Pointer(results)) + uintptr(i)*unsafe.Sizeof(Network{}),
			),
		)
		networks[i] = n
	}
	return networks
}

func activeCGO() string {
	return C.GoString(C.current_connection())
}

func connCGO(ssid, pass string, output chan string) bool {
	outputChan = output
	C.redirect_output()
	result := C.network_conn(C.CString(ssid), C.CString(pass)) == 0
	C.reset_output()

	return result
	// return C.conn(C.CString(ssid), C.CString(pass)) == 0
}
