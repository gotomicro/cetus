package x

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"

	"github.com/gotomicro/ego/core/elog"
	"go.uber.org/zap"
)

func S2I64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		elog.Warn("S2I64 fail", zap.Error(err))
		return 0
	}
	return i
}

type Int interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
}

func I2S[T Int](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func UI2S[T Uint](i T) string {
	return strconv.FormatUint(uint64(i), 10)
}

// Float64ToByte Float64转byte
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, bits)
	return data
}

// ByteToFloat64 byte转Float64
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	_ = binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

func Int64ToBytes(data int64) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.BigEndian, data)
	return bytebuf.Bytes()
}

func BytesToInt64(bys []byte) int64 {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	_ = binary.Read(bytebuff, binary.BigEndian, &data)
	return data
}
