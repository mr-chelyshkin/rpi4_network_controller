//go:build darwin
// +build darwin

package wifi

// Wifi controller object.
type Wifi struct{}

// NewWifi create Wifi object.
func NewWifi() (*Wifi, error) {
	return &Wifi{}, nil
}

// Scan wireless network.
func (w *Wifi) Scan() []*Network {
	return nil
}

// Conn to selected wireless network.
func (w *Wifi) Conn(_, _ string, _ chan string) bool {
	return false
}

// Active connection.
func (w *Wifi) Active() string {
	return "No connection"
}
