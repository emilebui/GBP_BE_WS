package test

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"testing"
)

var redisConn *redis.Client = nil

func getFakeRedis(t *testing.T) *redis.Client {
	if redisConn == nil {
		r, err := createFakeRedis()
		if err != nil {
			t.Fatalf("Failed to initiate fake redis!!!")
		}
		redisConn = r
	}

	return redisConn
}

func resetFakeRedis() {
	redisConn.FlushAll(context.Background())
}

func createFakeRedis() (*redis.Client, error) {
	s, err := miniredis.Run()
	if err != nil {
		return nil, err
	}
	r := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	return r, nil
}
