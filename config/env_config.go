package config

import (
	"fmt"
	"os"
	"sync"
)

type EnvConfig struct {
	ENV string
}

var (
	envInstance *EnvConfig
	envOnce     sync.Once
)

// ReadEnvConfig ensures the environment is read only once (singleton)
func ReadEnvConfig() *EnvConfig {
	envOnce.Do(func() {
		env := os.Getenv("ENV")
		if env == "" {
			env = "production" // Default to "production" if ENV is not set
		}

		envInstance = &EnvConfig{ENV: env}
		fmt.Println("Loaded ENV:", envInstance.ENV)
	})

	return envInstance
}
