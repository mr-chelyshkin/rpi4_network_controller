package wifi

import "C"
import (
	"fmt"
	"strconv"
)

func State(output chan string) string {
	return networkStateCGO(output)
}

func Conn(ssid, pass, country string, output chan string) bool {
	return networkConnCGO(ssid, pass, country, output)
}

type ScanOpts struct {
	SkipEmptySsid bool
}

func Scan(output chan string) []*Network {
	return networkScanCGO(output)
}

type Network struct {
	sSID    [33]C.char
	freq    float64
	quality int32
	level   int32
}

type ByLevelDesc []Network

func (a ByLevelDesc) Len() int           { return len(a) }
func (a ByLevelDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLevelDesc) Less(i, j int) bool { return a[i].level > a[j].level }

type BySsidDesc []Network

func (a BySsidDesc) Len() int           { return len(a) }
func (a BySsidDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySsidDesc) Less(i, j int) bool { return C.GoString(&a[i].sSID[0]) > C.GoString(&a[j].sSID[0]) }

// GetSSID return network ssid.
func (n Network) GetSSID() string {
	return C.GoString(&n.sSID[0])
}

// GetFreq return network freq.
func (n Network) GetFreq() string {
	return fmt.Sprintf("%.2f", n.freq/1e9)
}

// GetQuality return network quality.
func (n Network) GetQuality() string {
	return strconv.Itoa(int(n.quality))
}

// GetLevel return network level.
func (n Network) GetLevel() string {
	return strconv.Itoa(int(n.level))
}
