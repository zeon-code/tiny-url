package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type MetricConfiguration interface {
	Integration() (string, error)
	Environment() (string, error)
	Host() (string, error)
	Port() (int, error)
}

type DatadogConfiguration struct{}

func NewDatadogConfiguration() DatadogConfiguration {
	return DatadogConfiguration{}
}

func (c DatadogConfiguration) Environment() (string, error) {
	return c.get("ENV")
}

func (c DatadogConfiguration) Integration() (string, error) {
	return c.get("METRIC_INTEGRATION")
}

func (c DatadogConfiguration) Host() (string, error) {
	return c.get("METRIC_HOST")
}

func (c DatadogConfiguration) Port() (int, error) {
	env, err := c.get("METRIC_PORT")

	if err != nil {
		return 0, err
	}

	port, err := strconv.Atoi(env)

	if err != nil {
		return 0, errors.New("METRIC_PORT must be an interger value")
	}

	return port, nil
}

func (c DatadogConfiguration) get(env string) (string, error) {
	if value, exists := os.LookupEnv(env); exists {
		return value, nil
	}

	return "", fmt.Errorf("Missing required environment variable %s", env)
}
