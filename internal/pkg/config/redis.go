package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type RedisConfig struct{}

func (c RedisConfig) Driver() string {
	return "redis"
}

func (c RedisConfig) DSN() (string, error) {
	host, err := c.Host()

	if err != nil {
		return "", err
	}

	port, err := c.Port()

	if err != nil {
		return "", err
	}

	name, err := c.Name()

	if err != nil {
		return "", err
	}

	if password, err := c.Password(); err == nil {
		return fmt.Sprintf("redis://:%s@%s:%d/%d", password, host, port, name), nil
	}

	return fmt.Sprintf("redis://%s:%d/%d", host, port, name), nil
}

func (c RedisConfig) Host() (string, error) {
	return c.get("CACHE_HOST")
}

func (c RedisConfig) Port() (int, error) {
	env, err := c.get("CACHE_PORT")

	if err != nil {
		return 0, err
	}

	port, err := strconv.Atoi(env)

	if err != nil {
		return 0, errors.New("CACHE_PORT must be an interger value")
	}

	return port, nil
}

func (c RedisConfig) Password() (string, error) {
	return c.get("CACHE_PASSWORD")
}

func (c RedisConfig) Name() (int, error) {
	env, err := c.get("CACHE_NAME")

	if err != nil {
		return 0, err
	}

	name, err := strconv.Atoi(env)

	if err != nil {
		panic("CACHE_NAME must be an interger value")
	}

	return name, nil
}

func (c RedisConfig) get(env string) (string, error) {
	if value, exists := os.LookupEnv(env); exists {
		return value, nil
	}

	return "", fmt.Errorf("Missing required environment variable %s", env)
}
