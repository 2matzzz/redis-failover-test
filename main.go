package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisAddr     = ""
	redisPassword = ""
)

func main() {
	host := flag.String("h", "127.0.0.1", "Redis server IP address")
	port := flag.String("p", "6379", "Redis server port")
	flag.Parse()

	redisAddr = fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("Connecting to Redis at %s\n", redisAddr)

	go maintainConnectionPing()
	go newConnectionPing()
	select {}
}

func maintainConnectionPing() {
	ctx := context.Background()
	client := createRedisClient()

	for {
		pingStart := time.Now()
		pong, err := client.Ping(ctx).Result()
		if err != nil {
			log.Printf("maintainConnectionPing: Error - %v. Reconnecting...\n", err)
			client = reconnectRedis(client)
			continue
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

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		Password:    redisPassword,
		DB:          0,
		DialTimeout: 500 * time.Millisecond,
		ReadTimeout: 500 * time.Millisecond,
	})
}

func reconnectRedis(client *redis.Client) *redis.Client {
	client.Close()
	for {
		client = createRedisClient()
		_, err := client.Ping(context.Background()).Result()
		if err == nil {
			log.Println("Reconnected to Redis server successfully.")
			return client
		}
		log.Printf("Reconnect failed: %v. Retrying in 1 seconds...\n", err)
		time.Sleep(1 * time.Second)
	}
}

func newConnectionPing() {
	for {
		ctx := context.Background()
		client := createRedisClient()

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
