# XRateLimiter - Redis限流器

基于Redis的高性能限流器，使用滑动窗口算法实现精确限流控制。

## 功能特性

- ✅ 基于Redis的分布式限流
- ✅ 滑动窗口算法，精确控制请求频率
- ✅ 使用Lua脚本保证原子性操作
- ✅ 支持动态配置启用/禁用
- ✅ 自动清理过期数据，防止内存泄漏
- ✅ 配置灵活，支持自定义时间窗口和请求次数

## 安装

```bash
go get github.com/gotomicro/cetus/xratelimiter
```

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/gotomicro/cetus/xratelimiter"
)

func main() {
    // 创建Redis客户端
    redisClient := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    // 配置限流器
    config := xratelimiter.Config{
        Enable:   true,  // 启用限流
        Period:   10,    // 时间窗口：10秒
        MaxCount: 30,    // 最大请求次数：30次
    }

    // 创建限流器实例
    limiter, err := xratelimiter.NewRedisRateLimiter(redisClient, config)
    if err != nil {
        panic(err)
    }

    // 使用限流器
    userID := "user123"
    if limiter.Refused(userID) {
        fmt.Println("请求被限流")
    } else {
        fmt.Println("请求通过")
    }
}
```

### 配置说明

```go
type Config struct {
    Enable   bool `toml:"enable"`    // 是否启用限流，false时直接放行
    Period   int  `toml:"period"`    // 时间窗口（秒），默认10秒
    MaxCount int  `toml:"maxCount"`  // 时间窗口内最大请求次数，默认30次
}
```

## 详细示例

### 1. 基本限流示例

```go
func rateLimitExample() {
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // 每10秒最多5次请求
    limiter, _ := xratelimiter.NewRedisRateLimiter(redisClient, xratelimiter.Config{
        Enable:   true,
        Period:   10,
        MaxCount: 5,
    })

    key := "api_user123"
    
    // 模拟请求
    for i := 0; i < 10; i++ {
        if limiter.Refused(key) {
            fmt.Printf("请求 %d: 被限流 ❌\n", i+1)
        } else {
            fmt.Printf("请求 %d: 通过 ✅\n", i+1)
        }
        time.Sleep(time.Second)
    }
}
```

### 2. 动态配置示例

```go
func dynamicConfigExample() {
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // 初始配置：禁用限流
    config := xratelimiter.Config{
        Enable:   false,  // 初始禁用
        Period:   5,
        MaxCount: 3,
    }

    limiter, _ := xratelimiter.NewRedisRateLimiter(redisClient, config)

    // 此时所有请求都会通过
    fmt.Println("限流禁用状态:")
    for i := 0; i < 5; i++ {
        refused := limiter.Refused("test_key")
        fmt.Printf("请求 %d: %v\n", i+1, !refused)
    }
}
```

### 3. 不同场景的配置推荐

```go
// API接口限流 - 每分钟60次
apiLimiter, _ := xratelimiter.NewRedisRateLimiter(redisClient, xratelimiter.Config{
    Enable:   true,
    Period:   60,    // 60秒
    MaxCount: 60,    // 60次
})

// 登录限流 - 每小时5次
loginLimiter, _ := xratelimiter.NewRedisRateLimiter(redisClient, xratelimiter.Config{
    Enable:   true,
    Period:   3600,  // 1小时
    MaxCount: 5,     // 5次
})

// 短信发送限流 - 每天10次
smsLimiter, _ := xratelimiter.NewRedisRateLimiter(redisClient, xratelimiter.Config{
    Enable:   true,
    Period:   86400, // 24小时
    MaxCount: 10,    // 10次
})
```

### 4. Web框架集成示例

#### Gin中间件示例

```go
func RateLimitMiddleware(limiter xratelimiter.RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 使用用户ID或IP作为限流key
        key := c.GetHeader("User-ID")
        if key == "" {
            key = c.ClientIP()
        }

        if limiter.Refused(key) {
            c.JSON(429, gin.H{
                "error": "请求过于频繁，请稍后再试",
                "code":  429,
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// 使用中间件
func setupRouter() *gin.Engine {
    r := gin.Default()
    
    // 创建限流器
    redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
    limiter, _ := xratelimiter.NewRedisRateLimiter(redisClient, xratelimiter.Config{
        Enable:   true,
        Period:   60,
        MaxCount: 100,
    })

    // 应用限流中间件
    r.Use(RateLimitMiddleware(limiter))
    
    r.GET("/api/data", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    return r
}
```

## API文档

### 接口定义

```go
type RateLimiter interface {
    Refused(key string) bool
}
```

### 方法说明

#### `NewRedisRateLimiter(redisClient *redis.Client, limit Config) (RateLimiter, error)`

创建基于Redis的限流器实例。

**参数:**
- `redisClient`: Redis客户端实例
- `limit`: 限流配置

**返回:**
- `RateLimiter`: 限流器接口实例
- `error`: 错误信息

#### `Refused(key string) bool`

检查指定key是否应该被限流。

**参数:**
- `key`: 限流标识符（如用户ID、IP地址等）

**返回:**
- `bool`: true表示应该拒绝请求，false表示允许请求

## 配置说明

### 默认配置

如果配置值无效（<=0），将使用以下默认值：
- `Period`: 10秒
- `MaxCount`: 30次

### 配置文件示例（TOML）

```toml
[ratelimiter]
enable = true
period = 60        # 时间窗口：60秒
maxCount = 100     # 最大请求次数：100次
```

## 工作原理

### 滑动窗口算法

限流器使用滑动窗口算法，通过Redis的有序集合（ZSet）来实现：

1. **记录请求**: 每次请求时，将当前时间戳作为score和value存入ZSet
2. **清理过期**: 移除时间窗口之前的所有记录
3. **计数检查**: 统计当前窗口内的请求数量
4. **限流判断**: 如果超过限制则拒绝，否则记录本次请求

### Lua脚本

使用Lua脚本保证操作的原子性，避免并发问题：

```lua
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
return count+1
```

## 注意事项

1. **Redis连接**: 确保Redis服务可用且连接配置正确
2. **Key设计**: 合理设计限流key，避免冲突（建议添加业务前缀）
3. **时钟同步**: 分布式环境下确保各节点时钟同步
4. **内存管理**: 限流器会自动设置key的过期时间，防止内存泄漏
5. **错误处理**: 当Redis出现错误时，默认行为是拒绝请求（fail-safe）

## 性能特点

- **高性能**: 使用Lua脚本减少网络往返次数
- **内存优化**: 自动清理过期数据，避免内存泄漏
- **并发安全**: Lua脚本保证操作原子性
- **分布式**: 基于Redis，天然支持分布式限流

## 常见问题

### Q: 如何选择合适的时间窗口和请求次数？

A: 根据业务场景来定：
- **API接口**: 建议每分钟100-1000次
- **登录接口**: 建议每小时3-10次
- **短信发送**: 建议每天5-20次

### Q: 限流key如何设计？

A: 推荐格式：`业务模块:限流类型:标识符`
```go
key := fmt.Sprintf("api:user_request:%s", userID)
key := fmt.Sprintf("sms:send:%s", phoneNumber)
key := fmt.Sprintf("login:attempt:%s", clientIP)
```

### Q: Redis出现问题时会怎样？

A: 为了系统安全，当Redis出现错误时，限流器会默认拒绝请求。可以根据业务需求调整这个策略。

## 许可证

本项目使用 MIT 许可证。详情请参阅 [LICENSE](../LICENSE) 文件。
