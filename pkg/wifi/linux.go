//go:build linux

package wifi

/*
#cgo linux CFLAGS: -Dlinux
#cgo linux LDFLAGS: -liw
#include "linux_wifi.h"
*/
import "C"

import (
	"unsafe"
)

//export goSendToChannel
func goSendToChannel(s *C.char) {
	outputChan <- C.GoString(s)
}

var outputChan chan string

func networkStateCGO(output chan string) string {
	outputChan = output

	C.reset_output()
	defer C.reset_output()
	return C.GoString(C.network_state())
}

func networkConnCGO(ssid, pass, country string, output chan string) bool {
	outputChan = output

	C.redirect_output()
	defer C.reset_output()
	return C.network_conn(C.CString(ssid), C.CString(pass), C.CString(country)) == 0
}

func networkScanCGO(output chan string) []*Network {
	outputChan = output

	count := C.int(0)
	output <- "run C code"
	results := C.network_scan(&count)
	output <- "got result C code"
	networks := make([]*Network, count)

	C.redirect_output()
	defer C.reset_output()

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
