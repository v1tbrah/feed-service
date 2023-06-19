package config

import "os"

const (
	defaultHTTPServHost = "127.0.0.1"
	defaultHTTPServPort = "3939"
)
const (
	envNameHTTPServHost = "HTTP_HOST"
	envNameHTTPServPort = "HTTP_PORT"
)

type HTTP struct {
	Host string
	Port string
}

func newDefaultHTTPConfig() HTTP {
	return HTTP{
		Host: defaultHTTPServHost,
		Port: defaultHTTPServPort,
	}
}

func (c *HTTP) parseEnv() {
	envServHost := os.Getenv(envNameHTTPServHost)
	if envServHost != "" {
		c.Host = envServHost
	}

	envServPort := os.Getenv(envNameHTTPServPort)
	if envServPort != "" {
		c.Port = envServPort
	}
}
