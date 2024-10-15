package xgray

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"
	"github.com/pkg/errors"
)

func hashStringToNumber(inputString string) (uint64, error) {
	// 创建 SHA-256 哈希对象
	hasher := sha256.New()
	// 将字符串转换为字节数组并写入哈希对象
	hasher.Write([]byte(inputString))
	// 计算哈希值并转换为16进制字符串
	hashInHex := hex.EncodeToString(hasher.Sum(nil))
	// 将16进制字符串转换为整数
	var hashedNumber uint64
	_, err := fmt.Sscanf(hashInHex, "%016x", &hashedNumber)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to hash string to number: %s", inputString)
	}
	return hashedNumber, nil
}

// IsGray
//
//	@Description: 是否灰度, 传入计算灰度百分比的 key，根据 0-100 的比例进行灰度，hash 计算无法保证完全平均
//	@param key
//	@param percent
//	@return bool
func IsGray(key string, percent uint64) bool {
	if percent <= 0 {
		return false
	}
	hashedNumber, err := hashStringToNumber(key)
	if err != nil {
		elog.Error("xgray", elog.FieldMethod("IsGray"), l.S("key", key), l.E(err))
		return false
	}
	return hashedNumber%100 < percent
}
