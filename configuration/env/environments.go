package env

type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
	Testing     Environment = "testing"
)

var APP_ENVIRONMENT_KEY = "APP_ENV"
var APP_PORT_KEY = "APP_PORT"
var APP_HOST_KEY = "APP_HOST"
