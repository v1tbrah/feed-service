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
	GRPC        GRPC
	HTTP        HTTP
	Cache       Cache
	Kafka       Kafka
	RelationCli RelationCli
	PostCli     PostCli

	LogLvl zerolog.Level
}

func NewDefaultConfig() Config {
	return Config{
		GRPC:        newDefaultGRPCConfig(),
		HTTP:        newDefaultHTTPConfig(),
		Cache:       newDefaultCacheConfig(),
		Kafka:       newDefaultKafkaConfig(),
		RelationCli: newDefaultRelationCliConfig(),
		PostCli:     newDefaultPostCliConfig(),
		LogLvl:      defaultLogLvl,
	}
}

func (c *Config) ParseEnv() error {
	c.GRPC.parseEnv()

	c.HTTP.parseEnv()

	if err := c.Cache.parseEnv(); err != nil {
		return errors.Wrap(err, "cache config parse env")
	}

	c.Kafka.parseEnv()

	c.RelationCli.parseEnv()

	c.PostCli.parseEnv()

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
