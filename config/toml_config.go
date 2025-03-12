package config

import (
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

// Config struct to hold parsed configuration
type Config struct {
	Server struct {
		URL string `toml:"url"`
	} `toml:"server"`
}

var (
	instance *Config
	once     sync.Once
)

// ReadConfig loads the configuration only once (singleton)
func ReadConfig() (*Config, error) {
	var err error

	once.Do(func() {
		data, readErr := os.ReadFile("config.toml")
		if readErr != nil {
			err = readErr
			return
		}

		var cfg Config
		if unmarshalErr := toml.Unmarshal(data, &cfg); unmarshalErr != nil {
			err = unmarshalErr
			return
		}

		instance = &cfg
	})

	return instance, err
}
