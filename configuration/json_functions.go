package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"goserve/configuration/env"
)

func loadFromJSONFile(filename string) func(*Config) error {
	return func(config *Config) error {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return nil // Fichier optionnel
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("impossible de lire %s: %v", filename, err)
		}

		// Merger avec la configuration existante
		var fileConfig Config
		if err := json.Unmarshal(data, &fileConfig); err != nil {
			return fmt.Errorf("JSON invalide dans %s: %v", filename, err)
		}

		return mergeConfigs(config, &fileConfig)
	}
}

func loadFromEnvSpecificJSONFile() func(*Config) error {
	return func(config *Config) error {
		envValue := string(config.Server.Environment)
		if envVar := os.Getenv(env.APP_ENVIRONMENT_KEY); envVar != "" {
			envValue = strings.ToLower(envVar)
		}

		filename := fmt.Sprintf("config.%s.json", envValue)
		return loadFromJSONFile(filename)(config)
	}
}

func mergeConfigs(target, source *Config) error {
	if err := mergeStructs(&target.Server, &source.Server); err != nil {
		return fmt.Errorf("erreur merge server config: %v", err)
	}

	for key, value := range source.App {
		target.App[key] = value
	}

	for key, value := range source.Custom {
		target.Custom[key] = value
	}

	return nil
}

func mergeStructs(target, source interface{}) error {
	targetValue := reflect.ValueOf(target).Elem()
	sourceValue := reflect.ValueOf(source).Elem()

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceField := sourceValue.Field(i)
		targetField := targetValue.Field(i)

		if !sourceField.IsZero() && targetField.CanSet() {
			targetField.Set(sourceField)
		}
	}

	return nil
}
