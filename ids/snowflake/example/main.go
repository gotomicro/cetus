package main

import (
	"fmt"
)

const (
	twepoch        = int64(1483228800000)              //开始时间截 (2017-01-01)
	workerIdBits   = uint(10)                          //机器id所占的位数
	sequenceBits   = uint(12)                          //序列所占的位数
	workerIdMax    = int64(-1 ^ (-1 << workerIdBits))  //支持的最大机器id数量
	sequenceMask   = uint64(-1 ^ (-1 << sequenceBits)) //
	workerIdShift  = sequenceBits                      //机器id左移位数
	timestampShift = sequenceBits + workerIdBits       //时间戳左移位数
)

func main() {
	fmt.Println(workerIdMax)
	fmt.Println(sequenceMask)
	fmt.Println(1000 * 60 * 60 * 24 * 365 * 70)
	fmt.Println(1000*60*60*24*365*70 + twepoch)
	fmt.Println(1 << 41)
	fmt.Println((uint64(1397264040) << 22) / (1000 * 60 * 60 * 24 * 365))
	// fmt.Println(1024 % workerIdMax)
	// fmt.Println(1458 % workerIdMax)
	// fmt.Println(2457 % workerIdMax)
}
