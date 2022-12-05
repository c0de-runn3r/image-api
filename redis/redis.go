package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	*redis.Client
}

var ctx = context.Background()
var Redis = RedisDB{}

func SetupRedisClient() {
	Redis.Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (r *RedisDB) AddPair(id string, name string) {
	log.Printf("writing to redis [%s:%s]", id, name)
	err := r.Set(ctx, id, name, 0).Err()
	if err != nil {
		panic(err)
	}
}
func (r *RedisDB) GetValue(id string) string {
	val, err := r.Get(ctx, id).Result()
	log.Printf("getting from redis [%s:%s]", id, val)
	if err != nil {
		panic(err)
	}
	return val
}
