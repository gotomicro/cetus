package l

import (
	"go.uber.org/zap"
)

var (
	I   = zap.Int
	I64 = zap.Int64
	I32 = zap.Int32
	F64 = zap.Float64
	S   = zap.String
	B   = zap.Bool
	E   = zap.Error
	A   = zap.Any
)
