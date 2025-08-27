package xratelimiter

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	. "github.com/smartystreets/goconvey/convey"

)

func TestRateLimiter_Allowed(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{})
	Convey("redis rate limiter", t, func() {
		redisRateLimiter, err := NewRedisRateLimiter(redisClient, Config{
			Enable:   true,
			Period:   10,
			MaxCount: 1,
		})
		if err != nil {
			t.Fatalf("init redis rate limiter fail, err:%v", err)
		}

		key := "nEoooYe"
		for i := 0; i < 15; i++ {
			time.Sleep(time.Second)
			res := redisRateLimiter.Refused(key)
			if !res {
				fmt.Println("  ✅ 1 more requests allowed")
			} else {
				fmt.Printf("  ❌ Blocked \n")
			}
		}
	})
}
