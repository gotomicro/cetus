package kredis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Subscribes = make([]*Subscribe, 0)

type Subscribe struct {
	DbRedis  *redis.Client
	channels map[string]*redis.PubSub
}

// Init 全局 Subscribe 初始化
func Init(ctx context.Context, client *redis.Client, channels map[string]func(ctx context.Context, pubSub *redis.PubSub)) *Subscribe {
	s := &Subscribe{
		DbRedis:  client,
		channels: make(map[string]*redis.PubSub),
	}
	for channel, fc := range channels {
		ps := s.Sub(ctx, channel)
		s.channels[channel] = ps
		go fc(ctx, ps)
	}
	Subscribes = append(Subscribes, s)
	return s
}

// Close 消息订阅关闭
func Close() error {
	for _, s := range Subscribes {
		if s == nil {
			continue
		}
		for _, pubSub := range s.channels {
			_ = pubSub.Close()
		}
		_ = s.DbRedis.Close()
	}
	return nil
}

// Sub 订阅
func (s *Subscribe) Sub(ctx context.Context, channel string) *redis.PubSub {
	var err error
	pubSub := s.DbRedis.Subscribe(ctx, channel)
	_, err = pubSub.Receive(ctx)
	if err != nil {
		return nil
	}
	return pubSub
}

//// Pub 发布
//func (s *Subscribe) Pub(ctx context.Context, channel string, message []byte) error {
//	var err error
//	err = s.DbRedis.Publish(ctx, channel, message).Err()
//	if err != nil {
//		fmt.Printf("subscribe error: %s", err.Error())
//		return err
//	}
//	return nil
//}
