package controller

// ControllerOpts ...
type ControllerOpts func(*Controller)

// WithScanSkipEmptySSIDs do not show networks in scan result with empty SSID.
func WithScanSkipEmptySSIDs() ControllerOpts {
	return func(c *Controller) {
		c.scanSkipEmptySsid = true
	}
}
