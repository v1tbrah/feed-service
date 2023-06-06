package config

import "os"

const (
	defaultGRPCServHost = "127.0.0.1"
	defaultGRPCServPort = "3030"
)
const (
	envNameGRPCServHost = "GRPC_HOST"
	envNameGRPCServPort = "GRPC_PORT"
)

type GRPC struct {
	Host string
	Port string
}

func newDefaultGRPCConfig() GRPC {
	return GRPC{
		Host: defaultGRPCServHost,
		Port: defaultGRPCServPort,
	}
}

func (c *GRPC) parseEnv() {
	envServHost := os.Getenv(envNameGRPCServHost)
	if envServHost != "" {
		c.Host = envServHost
	}

	envServPort := os.Getenv(envNameGRPCServPort)
	if envServPort != "" {
		c.Port = envServPort
	}
}
