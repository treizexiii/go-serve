package configuration

import (
	"fmt"
	"goserve/configuration/utils"
	"log"
)

type ConfigurationSource struct {
	Filename string
	Priority int
	Load     func(*Config) error
}

type ConfigLoader struct {
	sources []ConfigurationSource
	config  *Config
}

func New() ConfigurationBuilder {
	return &ConfigLoader{
		sources: make([]ConfigurationSource, 0),
		config: &Config{
			App:    make(map[string]interface{}),
			Custom: make(map[string]interface{}),
		},
	}
}

func (cl *ConfigLoader) AddSource(source ConfigurationSource) ConfigurationBuilder {
	cl.sources = append(cl.sources, source)
	return cl
}

func (cl *ConfigLoader) Load(options ...ConfigOption) (Configuration, error) {
	cl.config.Server.setDefaults()
	cl.config.Server.loadFromEnv()

	for i := 0; i <= 10; i++ { // Priorité de 0 à 10
		for _, source := range cl.sources {
			if source.Priority == i {
				if err := source.Load(cl.config); err != nil {
					log.Printf("Error on loading %s: %v", source.Filename, err)
				} else {
					log.Printf("Loading configuration from: %s", source.Filename)
				}
			}
		}
	}

	for _, option := range options {
		option(cl.config)
	}

	cl.config.Server.validate()
	cl.config.Server.LogConfiguration()

	return cl.config, nil
}

func (cl *ConfigLoader) LoadConfig(options ...ConfigOption) (Configuration, error) {

	cl.AddSource(ConfigurationSource{
		Filename: "defaults",
		Priority: 0,
		Load:     func(c *Config) error { return nil }, // Déjà fait dans Load()
	})

	cl.AddSource(ConfigurationSource{
		Filename: "config.json",
		Priority: 1,
		Load:     loadFromJSONFile("config.json"),
	})

	cl.AddSource(ConfigurationSource{
		Filename: "config.{env}.json",
		Priority: 2,
		Load:     loadFromEnvSpecificJSONFile(),
	})

	cl.AddSource(ConfigurationSource{
		Filename: "environment variables",
		Priority: 3,
		Load:     loadFromEnvVars,
	})

	return cl.Load(options...)
}

func (cl *ConfigLoader) LoadConfigFromFile(filename string, options ...ConfigOption) (Configuration, error) {

	cl.AddSource(ConfigurationSource{
		Filename: "defaults",
		Priority: 0,
		Load:     func(c *Config) error { return nil },
	})

	cl.AddSource(ConfigurationSource{
		Filename: fmt.Sprintf("file: %s", filename),
		Priority: 1,
		Load:     loadFromJSONFile(filename),
	})

	cl.AddSource(ConfigurationSource{
		Filename: "environment variables",
		Priority: 2,
		Load:     loadFromEnvVars,
	})

	return cl.Load(options...)
}

func loadFromEnvVars(config *Config) error {
	return utils.LoadStructFromEnv(config)
}
