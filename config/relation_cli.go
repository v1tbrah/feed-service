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
	ServHost string
	ServPort string
}

func newDefaultRelationCliConfig() RelationCli {
	return RelationCli{
		ServHost: defaultRelationServiceClientHost,
		ServPort: defaultRelationServiceClientPort,
	}
}

func (u *RelationCli) parseEnv() {
	envServHost := os.Getenv(envNameRelationServiceClientHost)
	if envServHost != "" {
		u.ServHost = envServHost
	}

	envServPort := os.Getenv(envNameRelationServiceClientPort)
	if envServPort != "" {
		u.ServPort = envServPort
	}
}
