package health

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/pkg/errors"
)

const ModeName = "health"

type Config struct {
	Host string
	Port int

	logger *xlog.Logger
}

func DefaultConfig() *Config {
	return &Config{
		Host:   "127.0.0.1",
		Port:   8081,
		logger: xlog.JupiterLogger.With(xlog.FieldMod(ModeName)),
	}
}

func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil &&
		errors.Cause(err) != conf.ErrInvalidKey {
		config.logger.Panic("health parse config panic", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr))
	}
	return config
}

func (c *Config) WithLogger(logger *xlog.Logger) *Config {
	c.logger = logger
	return c
}

func (c *Config) WithHost(host string) *Config {
	c.Host = host
	return c
}

func (c *Config) WithPort(port int) *Config {
	c.Port = port
	return c
}
