package config

import "os"

const (
	defaultRelationServiceClientHost = "127.0.0.1"
	defaultRelationServiceClientPort = "4040"
)

const (
	envNameRelationServiceClientHost = "RELATION_SERVICE_CLIENT_HOST"
	envNameRelationServiceClientPort = "RELATION_SERVICE_CLIENT_PORT"
)

type RelationCli struct {
	Host string
	Port string
}

func newDefaultRelationCliConfig() RelationCli {
	return RelationCli{
		Host: defaultRelationServiceClientHost,
		Port: defaultRelationServiceClientPort,
	}
}

func (u *RelationCli) parseEnv() {
	envServHost := os.Getenv(envNameRelationServiceClientHost)
	if envServHost != "" {
		u.Host = envServHost
	}

	envServPort := os.Getenv(envNameRelationServiceClientPort)
	if envServPort != "" {
		u.Port = envServPort
	}
}
