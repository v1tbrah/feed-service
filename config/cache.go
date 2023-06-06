package config

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	defaultCacheHost          = "127.0.0.1"
	defaultCachePort          = "6379"
	defaultFeedDBIdx          = 7
	defaultUserToPostsMapName = "user_to_posts_map"
)
const (
	envNameCacheHost          = "CACHE_HOST"
	envNameCachePort          = "CACHE_PORT"
	envNameFeedDBIdx          = "FEED_DB_IDX"
	envNameUserToPostsMapName = "USER_TO_POSTS_MAP"
)

type Cache struct {
	Host               string
	Port               string
	FeedDBIdx          int
	UserToPostsMapName string
}

func newDefaultCacheConfig() Cache {
	return Cache{
		Host:               defaultCacheHost,
		Port:               defaultCachePort,
		FeedDBIdx:          defaultFeedDBIdx,
		UserToPostsMapName: defaultUserToPostsMapName,
	}
}

func (c *Cache) parseEnv() error {
	envHost := os.Getenv(envNameCacheHost)
	if envHost != "" {
		c.Host = envHost
	}

	envPort := os.Getenv(envNameCachePort)
	if envPort != "" {
		c.Port = envPort
	}

	envFeedDBIdxStr := os.Getenv(envNameFeedDBIdx)
	if envFeedDBIdxStr != "" {
		idx, err := strconv.Atoi(envFeedDBIdxStr)
		if err != nil {
			return errors.Wrapf(err, "parse feed db idx: %s", envFeedDBIdxStr)
		}
		c.FeedDBIdx = idx
	}

	envUserToPostsMapName := os.Getenv(envNameUserToPostsMapName)
	if envUserToPostsMapName != "" {
		c.UserToPostsMapName = envUserToPostsMapName
	}

	return nil
}
