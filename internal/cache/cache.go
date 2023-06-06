package cache

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"gitlab.com/pet-pr-social-network/feed-service/config"
)

type Cache struct {
	cli *redis.Client

	cfg config.Cache
}

func Init(cfg config.Cache) (*Cache, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: cfg.Host + ":" + cfg.Port,
		DB:   cfg.FeedDBIdx,
	})

	if err := cli.Ping(context.Background()).Err(); err != nil {
		fmt.Println("ВОТ ТУТА")
		return nil, errors.Wrapf(err, "cli.Ping: %s:%s", cfg.Host, cfg.Port)
	}

	return &Cache{
		cli: cli,
		cfg: cfg,
	}, nil
}

func (c *Cache) Close(ctx context.Context) (err error) {
	closeEnded := make(chan struct{})

	go func() {
		if err = c.cli.Close(); err != nil {
			err = errors.Wrap(err, "cli.Close")
			closeEnded <- struct{}{}
			return
		}

		closeEnded <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-closeEnded:
		return err
	}
}
