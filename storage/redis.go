package storage

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	*redis.Client
}

var ctx = context.Background()

func NewRedisClient(addr string) *RedisDB {
	client := *redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &RedisDB{&client}
}

func (r *RedisDB) AddPair(id string, name string) {
	log.Printf("writing to redis [%s:%s]", id, name)
	err := r.Set(ctx, id, name, 0).Err()
	if err != nil {
		log.Println("Error writing to redis", err)
	}
}

func (r *RedisDB) GetValue(id string) string {
	val, err := r.Get(ctx, id).Result()
	log.Printf("getting from redis [%s:%s]", id, val)
	if err != nil {
		log.Println("Error getting from redis", err)
	}
	return val
}
