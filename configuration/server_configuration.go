package configuration

import (
	"log"
	"os"
	"strings"
	"goserve/configuration/env"
)

type ServeurConfiguration struct {
	Environment env.Environment
	Port        int
	Host        string

	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

func (c *ServeurConfiguration) setDefaults() {
	c.Environment = env.Development
	c.Port = 8080
	c.Host = ""
	c.ReadTimeout = 15
	c.WriteTimeout = 15
	c.IdleTimeout = 60
}

func (c *ServeurConfiguration) loadFromEnv() {
	if envValue := os.Getenv(env.APP_ENVIRONMENT_KEY); envValue != "" {
		c.Environment = env.Environment(strings.ToLower(envValue))
	}
	if port := os.Getenv(env.APP_PORT_KEY); port != "" {
		c.Port = int(port[0])
	}
	if host := os.Getenv(env.APP_HOST_KEY); host != "" {
		c.Host = host
	}
	if rt := os.Getenv("READ_TIMEOUT"); rt != "" {
		c.ReadTimeout = int(rt[0])
	}
	if wt := os.Getenv("WRITE_TIMEOUT"); wt != "" {
		c.WriteTimeout = int(wt[0])
	}
	if it := os.Getenv("IDLE_TIMEOUT"); it != "" {
		c.IdleTimeout = int(it[0])
	}
}

func (c *ServeurConfiguration) validate() {
	validEnvs := map[env.Environment]bool{
		env.Development: true,
		env.Staging:     true,
		env.Production:  true,
		env.Testing:     true,
	}

	if !validEnvs[c.Environment] {
		panic("Invalid environment setting")
	}
}

func (c *ServeurConfiguration) LogConfiguration() {
	log.Println("Configuration loaded successfully")
	log.Printf("Environment: %s", c.Environment)
	log.Printf("Host: %s", c.Host)
	log.Printf("Port: %d", c.Port)
	log.Printf("Read Timeout: %d", c.ReadTimeout)
	log.Printf("Write Timeout: %d", c.WriteTimeout)
	log.Printf("Idle Timeout: %d", c.IdleTimeout)
}
