package config

import (
	"os"
)

const (
	defaultKafkaEnable        = true
	defaultKafkaHost          = "127.0.0.1"
	defaultKafkaPort          = "9092"
	defaultTopicPostCreated   = "post_created"
	defaultTopicPostDeleted   = "post_deleted"
	defaultTopicFriendAdded   = "friend_added"
	defaultTopicFriendRemoved = "friend_removed"
)
const (
	envNameKafkaEnable        = "KAFKA_ENABLE"
	envNameKafkaHost          = "KAFKA_HOST"
	envNameKafkaPort          = "KAFKA_PORT"
	envNameTopicPostCreated   = "TOPIC_POST_CREATED"
	envNameTopicPostDeleted   = "TOPIC_POST_DELETED"
	envNameTopicFriendAdded   = "TOPIC_FRIEND_ADDED"
	envNameTopicFriendRemoved = "TOPIC_FRIEND_REMOVED"
)

type KafkaConfig struct {
	Enable bool

	Host string
	Port string

	TopicPostCreated string
	TopicPostDeleted string

	TopicFriendAdded   string
	TopicFriendRemoved string
}

func newDefaultKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Enable:             defaultKafkaEnable,
		Host:               defaultKafkaHost,
		Port:               defaultKafkaPort,
		TopicPostCreated:   defaultTopicPostCreated,
		TopicPostDeleted:   defaultTopicPostDeleted,
		TopicFriendAdded:   defaultTopicFriendAdded,
		TopicFriendRemoved: defaultTopicFriendRemoved,
	}
}

func (c *KafkaConfig) parseEnv() {
	envKafkaEnable := os.Getenv(envNameKafkaEnable)
	if envKafkaEnable != "" {
		c.Enable = envKafkaEnable == "true"
	}

	envKafkaHost := os.Getenv(envNameKafkaHost)
	if envKafkaHost != "" {
		c.Host = envKafkaHost
	}

	envKafkaPort := os.Getenv(envNameKafkaPort)
	if envKafkaPort != "" {
		c.Port = envKafkaPort
	}

	envTopicPostCreated := os.Getenv(envNameTopicPostCreated)
	if envTopicPostCreated != "" {
		c.TopicPostCreated = envTopicPostCreated
	}

	envTopicPostDeleted := os.Getenv(envNameTopicPostDeleted)
	if envTopicPostDeleted != "" {
		c.TopicPostDeleted = envTopicPostDeleted
	}

	envTopicFriendAdded := os.Getenv(envNameTopicFriendAdded)
	if envTopicFriendAdded != "" {
		c.TopicFriendAdded = envTopicFriendAdded
	}

	envTopicFriendRemoved := os.Getenv(envNameTopicFriendRemoved)
	if envTopicFriendRemoved != "" {
		c.TopicFriendRemoved = envTopicFriendRemoved
	}
}
