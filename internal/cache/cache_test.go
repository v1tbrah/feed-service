//go:build with_db

package cache

import (
	"context"
	"testing"

	"github.com/rs/zerolog"

	"gitlab.com/pet-pr-social-network/feed-service/config"
)

func tHelperInitEmptyCache(t *testing.T) *Cache {
	cfg := config.NewDefaultConfig()
	zerolog.SetGlobalLevel(cfg.LogLvl)

	if err := cfg.ParseEnv(); err != nil {
		t.Fatalf("config.ParseEnv: %v", err)
	}
	zerolog.SetGlobalLevel(cfg.LogLvl)

	c, err := Init(cfg.CacheConfig)
	if err != nil {
		t.Fatalf("init cache: %v", err)
	}

	if err = c.cli.FlushDB(context.Background()).Err(); err != nil {
		t.Fatalf("flush db: %v", err)
	}

	return c
}
