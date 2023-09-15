//go:build linux

package wifi

/*
#cgo linux CFLAGS: -Dlinux
#cgo linux LDFLAGS: -liw

#include "linux_wifi.h"
*/
import "C"

var outputCh chan string

func goSendToChannel(s *C.char) {
	outputCh <- C.GoString(s)
}

func scanCGO() []*Network {
	networks := make([]*Network, count)
}

// Wifi object
type Wifi struct{}

// NewWifi init and return Wifi object
func NewWifi() (*Wifi, error) {
	return &Wifi{}, nil
}

func (w *Wifi) Scan() []*Network {

}
