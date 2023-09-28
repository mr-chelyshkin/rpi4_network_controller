//go:build linux

package wifi

/*
#cgo linux CFLAGS: -Dlinux
#cgo linux LDFLAGS: -liw
#include "linux_wifi.h"
*/
import "C"

//export goSendToChannel
func goSendToChannel(s *C.char) {
	outputChan <- C.GoString(s)
}

var outputChan chan string

func networkStateCGO(output chan string) string {
	outputChan = output

	C.redirect_output()
	defer C.reset_output()
	return C.GoString(C.current_connection())
}

func networkConnCGO(ssid, pass string, output chan string) bool {
	outputChan = output

	C.redirect_output()
	defer C.reset_output()
	return C.network_conn(C.CString(ssid), C.CString(pass)) == 0
}

func networkScanCGO(output chan string) []*Network {
	outputChan = output
	count := C.int(0)

	results := C.network_scan(&count)
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