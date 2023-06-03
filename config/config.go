package config

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	defaultLogLvl = zerolog.InfoLevel
	envNameLogLvl = "LOG_LVL"
)

type Config struct {
	GRPCConfig        GRPCConfig
	CacheConfig       CacheConfig
	KafkaConfig       KafkaConfig
	RelationCliConfig RelationCli
	PostCliConfig     PostCli

	LogLvl zerolog.Level
}

func NewDefaultConfig() Config {
	return Config{
		GRPCConfig:        newDefaultGRPCConfig(),
		CacheConfig:       newDefaultCacheConfig(),
		KafkaConfig:       newDefaultKafkaConfig(),
		RelationCliConfig: newDefaultRelationCliConfig(),
		PostCliConfig:     newDefaultPostCliConfig(),
		LogLvl:            defaultLogLvl,
	}
}

func (c *Config) ParseEnv() error {
	c.GRPCConfig.parseEnv()

	if err := c.CacheConfig.parseEnv(); err != nil {
		return errors.Wrap(err, "cache config parse env")
	}

	c.KafkaConfig.parseEnv()

	c.RelationCliConfig.parseEnv()

	c.PostCliConfig.parseEnv()

	if err := c.parseEnvLogLvl(); err != nil {
		return err
	}

	return nil
}

func (c *Config) parseEnvLogLvl() error {
	envLogLvl := os.Getenv(envNameLogLvl)
	if envLogLvl != "" {
		logLevel, err := zerolog.ParseLevel(envLogLvl)
		if err != nil {
			return errors.Wrapf(err, "parse log lvl: %s", envLogLvl)
		}
		c.LogLvl = logLevel
	}
	return nil
}
