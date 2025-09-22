package configuration

import "goserve/configuration/env"

type ConfigurationBuilder interface {
	// Add source configuration
	// File source could be JSON, YAML, etc.
	// The priority determines the order of loading (0 = highest, 10 = lowest)
	AddSource(source ConfigurationSource) ConfigurationBuilder
	// Load and return the final configuration
	Load(options ...ConfigOption) (Configuration, error)
	// Convenience method to load configuration from default files
	LoadConfig(options ...ConfigOption) (Configuration, error)
	// Convenience method to load configuration from a specific file
	LoadConfigFromFile(path string, options ...ConfigOption) (Configuration, error)
}

type Configuration interface {
	// Retrieve the server address
	GetAddress() string
	// Retrieve the server port
	GetPort() int
	// Retrieve the application environment
	GetEnvironment() env.Environment
	// Retrieve the server read timeout in seconds
	GetReadTimeout() int
	// Retrieve the server write timeout in seconds
	GetWriteTimeout() int
	// Retrieve the server idle timeout in seconds
	GetIdleTimeout() int

	// Check if the environment is development
	IsDevelopment() bool
	// Check if the environment is staging
	IsStaging() bool
	// Check if the environment is production
	IsProduction() bool
	// Check if the environment is testing
	IsTesting() bool
}
