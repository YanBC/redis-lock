package main

import (
	"context"
	"sync"
	"time"

	"github.com/YanBC/redis-lock/rlock"
	"github.com/go-redis/redis/v8"
)

func worker(num int) {
	lock_name := "lock:counter"
	redis_addr := "172.17.0.5:6379"
	redis_passwd := ""
	redis_db := 0
	expiration_time := 5 * time.Second
	lock := rlock.NewLock(lock_name, redis_addr, redis_passwd, redis_db, expiration_time)
	cache_name := "cache:counter"

	cli := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_passwd,
		DB:       redis_db,
	})
	ctx := context.Background()
	for !lock.Acquire(ctx) {
	}
	defer lock.Release(ctx)

	current, err := cli.Get(ctx, cache_name).Int()
	if err != nil {
		panic(err)
	}
	new := current + num
	cli.Set(ctx, cache_name, new, 0)
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			worker(num)
		}(i)
	}

	wg.Wait()
}
