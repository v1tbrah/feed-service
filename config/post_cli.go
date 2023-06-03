package config

import "os"

const (
	defaultPostServiceClientHost = "127.0.0.1"
	defaultPostServiceClientPort = "5050"
)

const (
	envNamePostServiceClientHost = "RELATION_SERVICE_CLIENT_HOST"
	envNamePostServiceClientPort = "RELATION_SERVICE_CLIENT_PORT"
)

type PostCli struct {
	ServHost string
	ServPort string
}

func newDefaultPostCliConfig() PostCli {
	return PostCli{
		ServHost: defaultPostServiceClientHost,
		ServPort: defaultPostServiceClientPort,
	}
}

func (u *PostCli) parseEnv() {
	envServHost := os.Getenv(envNamePostServiceClientHost)
	if envServHost != "" {
		u.ServHost = envServHost
	}

	envServPort := os.Getenv(envNamePostServiceClientPort)
	if envServPort != "" {
		u.ServPort = envServPort
	}
}
