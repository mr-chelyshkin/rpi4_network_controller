package wifi

/*
#cgo LDFLAGS: -liw
#include "wifi.h"
*/
import "C"
import (
	"fmt"
	"sort"
	"strconv"
	"unsafe"
)

// Wifi ...
type Wifi struct{}

// NewWifi ...
func NewWifi() (Wifi, error) {
	return Wifi{}, nil
}

// Scan ...
func (w *Wifi) Scan() []*network {
	uniqueMap := make(map[string]struct{})
	var uniqueRes []*network

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
func (w *Wifi) Conn(ssid, password string) bool {
	return connCGO(ssid, password)
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

type network struct {
	sSID    [33]C.char
	freq    float64
	quality int32
	level   int32
}

type byLevelDesc []*network

func (a byLevelDesc) Len() int           { return len(a) }
func (a byLevelDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byLevelDesc) Less(i, j int) bool { return a[i].level > a[j].level }

// GetSSID ...
func (n *network) GetSSID() string {
	return C.GoStringN(&n.sSID[0], 32)
}

// GetFreq ...
func (n *network) GetFreq() string {
	return fmt.Sprintf("%.2f", n.freq/1e9)
}

// GetQuality ...
func (n *network) GetQuality() string {
	return strconv.Itoa(int(n.quality))
}

// GetLevel ...
func (n *network) GetLevel() string {
	return strconv.Itoa(int(n.level))
}

func scanCGO() []*network {
	count := C.int(0)
	results := C.scan(&count)
	networks := make([]*network, count)

	for i := 0; i < int(count); i++ {
		n := (*network)(
			unsafe.Pointer(
				uintptr(unsafe.Pointer(results)) + uintptr(i)*unsafe.Sizeof(network{}),
			),
		)
		networks[i] = n
	}
	return networks
}

func activeCGO() string {
	return C.GoString(C.active())
}

func connCGO(ssid, pass string) bool {
	return C.conn(C.CString(ssid), C.CString(pass)) == 0
}
