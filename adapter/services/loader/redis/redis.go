package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vladimir3322/stonent_go/config"
	"time"
)

var ctx = context.Background()
var client = redis.NewClient(&redis.Options{
	Addr: config.RedisUrl,
})

func PushEvent(data []byte) {
	client.RPush(ctx, config.RedisJobQueue, data)
}

func ConsumeEvents() {
	var ctx = context.Background()
	var client = redis.NewClient(&redis.Options{
		Addr: config.RedisUrl,
	})

	go func() {
		for {
			result, err := client.BLPop(ctx, 0*time.Second, config.RedisJobQueue).Result()

			if err != nil {
				fmt.Println("Error with event receiving:", err)
				continue
			}

			imageSource := result[1]

			fmt.Println("Image source size from redis: ", len(imageSource))
		}
	}()
}
