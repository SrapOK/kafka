package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func New(opt *redis.Options) (*redis.Client, error) {
	rdb := redis.NewClient(opt)
	status := rdb.Ping(context.Background())

	if err := status.Err(); err != nil {
		return nil, fmt.Errorf("failed to connect: %s", status.Err())
	}

	return rdb, nil
}
