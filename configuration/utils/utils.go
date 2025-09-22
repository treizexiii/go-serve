package utils

import (
	"goserve/configuration/env"
	"log"
	"os"
	"reflect"
	"strings"
)

func LoadStructFromEnv(target interface{}) error {
	value := reflect.ValueOf(target).Elem()
	typ := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typ.Field(i)

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}

		if err := setFieldFromString(field, envValue); err != nil {
			log.Printf("⚠️ Error setting field %s from env var %s: %v", fieldType.Name, envTag, err)
		}
	}

	return nil
}

func setFieldFromString(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		b := parseBool(value)
		field.SetBool(b)

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			parts := strings.Split(value, ",")
			slice := make([]string, len(parts))
			for i, part := range parts {
				slice[i] = strings.TrimSpace(part)
			}
			field.Set(reflect.ValueOf(slice))
		}

	case reflect.Struct:
		// Special case for Duration
		// if field.Type().Name() == "Duration" {
		// 	if d, err := time.ParseDuration(value); err == nil {
		// 		field.Set(reflect.ValueOf(Duration{d}))
		// 	}
		// }
		// Special case for Environment
		if field.Type().Name() == "Environment" {
			field.Set(reflect.ValueOf(env.Environment(strings.ToLower(value))))
		}
	}

	return nil
}

func parseBool(value string) bool {
	trueValues := []string{"1", "t", "true", "yes", "y"}
	value = strings.ToLower(strings.TrimSpace(value))
	for _, v := range trueValues {
		if value == v {
			return true
		}
	}
	return false
}
