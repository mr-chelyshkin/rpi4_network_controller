//go:build darwin

package wifi

func networkStateCGO(output chan string) string {
	return ""
}

func networkConnCGO(ssid, pass, country string, output chan string) bool {
	return false
}

func networkScanCGO(output chan string) []*Network {
	return nil
}
