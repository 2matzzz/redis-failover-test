package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisAddr = "localhost:6379"
var redisPassword = ""

func main() {
	go maintainConnectionPing()
	go newConnectionPing()
	select {}
}

func maintainConnectionPing() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
	defer client.Close()

	for {
		pingStart := time.Now()
		pong, err := client.Ping(ctx).Result()
		if err != nil {
			log.Printf("maintainConnectionPing: Error - %v\n", err)
			return
		}
		pingDuration := time.Since(pingStart)
		if pingDuration > 1*time.Second {
			log.Printf("maintainConnectionPing: Pong response delayed - took %v\n", pingDuration)
		} else {
			log.Printf("maintainConnectionPing: %s\n", pong)
		}
		time.Sleep(1 * time.Second)
	}
}

func newConnectionPing() {
	for {
		ctx := context.Background()
		client := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       0,
		})

		pong, err := client.Ping(ctx).Result()
		if err != nil {
			log.Printf("newConnectionPing: Error - %v\n", err)
		} else {
			log.Printf("newConnectionPing: %s\n", pong)
		}

		client.Close()
		time.Sleep(1 * time.Second)
	}
}
