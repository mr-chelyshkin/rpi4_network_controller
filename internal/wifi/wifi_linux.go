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

var outputCh chan string

func goSendToChannel(s *C.char) {
	outputCh <- C.GoString(s)
}

func networkActiveCGO() string {
	return C.GoString(C.current_connection())
}

func networkConnCGO(ssid, pass string, output chan string) bool {
	outputCh = output

	C.redirect_output()
	defer C.reset_output()
	return C.network_conn(C.CString(ssid), C.CString(pass)) == 0
}

func networkScanCGO() []*Network {
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

// Wifi object.
type Wifi struct{}

// NewWifi return Wifi object.
func NewWifi() (*Wifi, error) {
	return &Wifi{}, nil
}

// Active return current wifi connection.
func (w *Wifi) Active() string {
	return networkActiveCGO()
}

// Conn try to connect to selected network.
func (w *Wifi) Conn(ssid, password string, output chan string) bool {
	return networkConnCGO(ssid, password, output)
}

// Scan make scan wireless networks.
func (w *Wifi) Scan() []*Network {
	var results []*Network

	m := make(map[string]struct{})
	for _, net := range networkScanCGO() {
		if C.GoString(&net.sSID[0]) == "" {
			continue
		}

		key := C.GoString(&net.sSID[0]) + "_" + fmt.Sprintf("%f", net.freq)
		if _, found := m[key]; !found {
			results = append(results, net)
			m[key] = struct{}{}
		}
	}
	sort.Sort(byLevelDesc(results))
	return results
}
