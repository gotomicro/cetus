package l

import (
	"go.uber.org/zap"
)

var (
	I   = zap.Int
	I8  = zap.Int8
	I16 = zap.Int16
	I32 = zap.Int32
	I64 = zap.Int64
	U   = zap.Uint
	U8  = zap.Uint8
	U16 = zap.Uint16
	U32 = zap.Uint32
	U64 = zap.Uint64
	F64 = zap.Float64
	S   = zap.String
	B   = zap.Bool
	E   = zap.Error
	A   = zap.Any
	D   = zap.Duration
	T   = zap.Time
)
