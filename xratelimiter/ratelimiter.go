package xratelimiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"
)

type Config struct {
	Enable   bool `toml:"enable"`
	Period   int  `toml:"period"`
	MaxCount int  `toml:"maxCount"`
}

type RateLimiter interface {
	Refused(key string) bool
}

// -- 记录行为
// ; -- value 和 score 都使用纳秒时间戳，即ARGV[1]
// ; -- 移除时间窗口之前的行为记录，剩下的都是时间窗口内的
// ; -- 获取窗口内的行为数量
// ; -- 设置 zset 过期时间，避免冷用户持续占用内存
// -- 返回窗口内行为数量
var counterLuaScript = `
	-- 移除时间窗口之前的记录
	redis.pcall("zremrangebyscore", KEYS[1], 0, ARGV[2])
	-- 获取当前窗口内的记录数量
	local count = redis.call("zcard", KEYS[1])
	-- 如果已达到限制，直接返回当前数量（不增加）
	if count >= tonumber(ARGV[4]) then
		redis.pcall("expire", KEYS[1], ARGV[3])
		return count+1
	end
	-- 添加新记录
	redis.pcall("zadd", KEYS[1], ARGV[1], ARGV[1])
	-- 设置过期时间
	redis.pcall("expire", KEYS[1], ARGV[3])
	-- 返回更新后的数量
	return count+1`

// period单位为秒
type rateLimiter struct {
	evalSha     string
	period      int
	maxCount    int
	redisClient *redis.Client
	enable      bool
}

func NewRedisRateLimiter(redisClient *redis.Client, limit Config) (RateLimiter, error) {
	evalSha, err := redisClient.ScriptLoad(context.Background(), counterLuaScript).Result()
	if err != nil {
		return nil, err
	}

	// 配置默认值，10秒 30次请求
	period := limit.Period
	if period <= 0 {
		period = 10
	}
	maxCount := limit.MaxCount
	if maxCount <= 0 {
		maxCount = 30
	}

	return &rateLimiter{
		period:      period,
		maxCount:    maxCount,
		evalSha:     evalSha,
		redisClient: redisClient,
		enable:      limit.Enable,
	}, nil
}

func (r *rateLimiter) Refused(key string) bool {
	if !r.enable {
		return false
	}
	now := time.Now().UnixNano()
	beforeTime := now - int64(r.period*1000000000)
	res, err := r.redisClient.EvalSha(context.Background(), r.evalSha, []string{key}, now, beforeTime, r.period, r.maxCount).Result()
	if err != nil {
		elog.Warn("ratelimiter_err redis evalsha err", l.S("key", key), l.E(err))
		return true
	}
	// 如果当前计数大于等于最大限制，则拒绝请求
	if res.(int64) > int64(r.maxCount) {
		return true
	}
	return false
}
