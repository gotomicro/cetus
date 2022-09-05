package bucket

import (
	"sync"
	"time"
)

// Bucket 令牌桶
type Bucket struct {
	mu        sync.Mutex
	tokens    int64
	timestamp int64
	cap       int64
	rate      float32
}

// NewBucket 初始化令牌桶，速率为tokens/ms
func NewBucket(cap int64, rate float32) *Bucket {
	return &Bucket{
		cap:  cap,
		rate: rate,
	}
}

// Update 更新桶容量及速率
func (b *Bucket) Update(cap int64, rate float32) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.cap = cap
	b.rate = rate
}

// GetToken 获取报文令牌
func (b *Bucket) GetToken() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now().UnixNano() / 1e6

	tokens := int64(float32(now-b.timestamp) * b.rate)

	if tokens > 0 {
		b.timestamp = now
		b.tokens = b.tokens + tokens
	}

	if b.tokens > b.cap {
		b.tokens = b.cap
	}

	if b.tokens < 1 {
		return false
	}

	b.tokens--

	return true
}
