package config

import "os"

const (
	defaultPostServiceClientHost = "127.0.0.1"
	defaultPostServiceClientPort = "5050"
)

const (
	envNamePostServiceClientHost = "POST_SERVICE_CLIENT_HOST"
	envNamePostServiceClientPort = "POST_SERVICE_CLIENT_PORT"
)

type PostCli struct {
	Host string
	Port string
}

func newDefaultPostCliConfig() PostCli {
	return PostCli{
		Host: defaultPostServiceClientHost,
		Port: defaultPostServiceClientPort,
	}
}

func (u *PostCli) parseEnv() {
	envServHost := os.Getenv(envNamePostServiceClientHost)
	if envServHost != "" {
		u.Host = envServHost
	}

	envServPort := os.Getenv(envNamePostServiceClientPort)
	if envServPort != "" {
		u.Port = envServPort
	}
}
