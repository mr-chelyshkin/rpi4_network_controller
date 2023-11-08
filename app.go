package rpi4_network_controller

const (
	ScanTimeoutSec = 4
	ScanTickGlobal = 4
	Version        = "0.0.1"

	CtxKeyHotkeys  = "hotkeys"
	CtxKeyCurConn  = "conn"
	CtxKeyOutputCh = "output"
	CtxKeyUserName = "username"
	CtxKeyUserUid  = "useruid"
)

var (
	UserName = ""
	UserPerm = ""
)
