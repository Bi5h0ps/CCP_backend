package main

import (
	"CCP_backend/backend"
	"github.com/go-redis/redis"
)

func main() {
	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()
	server := backend.NewGin(client)
	server.Init()
	server.Start(":9090")
}
