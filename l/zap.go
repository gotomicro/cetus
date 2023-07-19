package l

import (
	"go.uber.org/zap"
)

var (
	I    = zap.Int
	I32  = zap.Int32
	UI32 = zap.Uint32
	I64  = zap.Int64
	UI64 = zap.Uint64
	F64  = zap.Float64
	S    = zap.String
	B    = zap.Bool
	E    = zap.Error
	A    = zap.Any
)
