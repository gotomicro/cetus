package idgen

import (
	"sync"
)

const (
	// InvalidID 非法ID
	InvalidID = 0
)

// IDBucket ID桶
type IDBucket struct {
	used          int // 已经使用的连接个数
	curIndex      int
	idPool        []uint32
	idRecyclePool []uint32
	idRecycleBits []uint64
	mu            sync.Mutex
	max           uint32
}

// New 创建id桶
func New(max uint32) *IDBucket {
	obj := &IDBucket{
		curIndex:      0,
		idPool:        make([]uint32, max),
		idRecyclePool: make([]uint32, 0, max),
		idRecycleBits: make([]uint64, 0, max/64+1),
	}

	var id uint32 = 1
	for i := 0; i < int(max); i++ {
		obj.idPool[i] = id
		id++
	}
	return obj
}

// Get 获取可用ID
func (c *IDBucket) Get() uint32 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.curIndex == len(c.idPool) {
		c.idPool, c.idRecyclePool = c.idRecyclePool, c.idPool[:0]
		c.curIndex = 0
		for i := 0; i < len(c.idRecycleBits); i++ {
			c.idRecycleBits[i] = 0
		}
	}

	id := c.idPool[c.curIndex]
	c.curIndex++
	c.used++
	return id
}

// Release 释放可用ID
func (c *IDBucket) Release(id uint32) bool {
	if id == InvalidID || id > c.max {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.idRecycleBits[id>>6]&(1<<(id&(64-1))) != 0 {
		return false
	}

	c.idRecyclePool = append(c.idRecyclePool, id)
	c.idRecycleBits[id>>6] |= 1 << (id & (64 - 1))

	if c.used > 0 {
		c.used--
	}

	return true
}
