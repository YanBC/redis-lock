package rlock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Lock struct {
	name       string
	cli        *redis.Client
	expiration time.Duration
	luuid      string
}

func (lock *Lock) Acquire(ctx context.Context) bool {
	success, err := lock.cli.SetNX(ctx, lock.name, lock.luuid, lock.expiration).Result()
	if err != nil {
		panic(err)
	}
	return success
}

func (lock *Lock) Release(ctx context.Context) bool {
	script := `if redis.call("get",KEYS[1]) == ARGV[1]
	then
		return redis.call("del",KEYS[1])
	else
		return 0
	end`
	vals, err := lock.cli.Eval(ctx, script, []string{lock.name}, lock.luuid).Result()
	if err != nil {
		panic(err)
	}
	num, isOK := vals.(int64)
	if !isOK {
		fmt.Println("wrong type")
	}
	return num != 0
}

func NewLock(name string, addr string, password string, db int, exp_time time.Duration) *Lock {
	ctx := context.Background()

	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return &Lock{
		name:       name,
		cli:        cli,
		expiration: exp_time,
		luuid:      uuid.New().String(),
	}
}
