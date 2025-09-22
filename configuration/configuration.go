package configuration

import "goserve/configuration/env"

type Config struct {
	Server ServeurConfiguration
	App    map[string]interface{} `jon:"app" yaml:"app"`
	Custom map[string]interface{} `json:"custom" yaml:"custom"`
}

type ConfigOption func(*Config)

func (c *Config) GetAddress() string {
	return c.Server.Host
}

func (c *Config) GetEnvironment() env.Environment {
	return c.Server.Environment
}

func (c *Config) GetIdleTimeout() int {
	return c.Server.IdleTimeout
}

func (c *Config) GetPort() int {
	return c.Server.Port
}

func (c *Config) GetReadTimeout() int {
	return c.Server.ReadTimeout
}

func (c *Config) GetWriteTimeout() int {
	return c.Server.WriteTimeout
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == env.Development
}

func (c *Config) IsStaging() bool {
	return c.Server.Environment == env.Staging
}

func (c *Config) IsProduction() bool {
	return c.Server.Environment == env.Production
}

func (c *Config) IsTesting() bool {
	return c.Server.Environment == env.Testing
}
