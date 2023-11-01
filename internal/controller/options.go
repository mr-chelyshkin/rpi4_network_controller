package controller

// ControllerOpts ...
type ControllerOpts func(*Controller)

// WithScanSkipEmptySSIDs do not show networks in scan result with empty SSID.
func WithScanSkipEmptySSIDs() ControllerOpts {
	return func(c *Controller) {
		c.scanSkipEmptySsid = true
	}
}

// WithScanSortByLevel sort scan results by signal level, desc.
func WithScanSortByLevel() ControllerOpts {
	return func(c *Controller) {
		c.scanSortBySignalLvl = true
	}
}

// WithScanSortByName sort scan results by network SSID, desc.
func WithScanSortByName() ControllerOpts {
	return func(c *Controller) {
		c.scanSortBySsidName = true
	}
}
