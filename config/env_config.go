package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
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
		// Load .env file if exists
		if err := godotenv.Load(); err != nil {
			fmt.Println("Warning: No .env file found, using system environment variables")
		}

		env := os.Getenv("ENV")
		if env == "" {
			env = "production" // Default is now "production"
		}

		envInstance = &EnvConfig{ENV: env}
		fmt.Println("Loaded ENV:", envInstance.ENV)
	})

	return envInstance
}
