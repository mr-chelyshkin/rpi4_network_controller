package wifi

import "C"
import (
	"fmt"
	"strconv"
)

type Network struct {
	sSID    [33]C.char
	freq    float64
	quality int32
	level   int32
}

type byLevelDesc []*Network

func (a byLevelDesc) Len() int           { return len(a) }
func (a byLevelDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byLevelDesc) Less(i, j int) bool { return a[i].level > a[j].level }

// GetSSID return network ssid.
func (n *Network) GetSSID() string {
	return C.GoStringN(&n.sSID[0], 32)
}

// GetFreq return network freq.
func (n *Network) GetFreq() string {
	return fmt.Sprintf("%.2f", n.freq/1e9)
}

// GetQuality return network quality.
func (n *Network) GetQuality() string {
	return strconv.Itoa(int(n.quality))
}

// GetLevel return network level.
func (n *Network) GetLevel() string {
	return strconv.Itoa(int(n.level))
}
