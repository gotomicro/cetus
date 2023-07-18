package xsocket

const (
	FcOpen    = 0
	FcClose   = 1
	FcPing    = 2
	FcPong    = 3
	FcMessage = 4
	FcUpgrade = 5
	FcNoop    = 6
)

const (
	ScEmpty       = -1
	ScConnect     = 0
	ScDisconnect  = 1
	ScEvent       = 2
	ScAck         = 3
	ScError       = 4
	ScBinaryEvent = 5
	ScBinaryAck   = 6
)

const (
	EventMessage = "message"
)
