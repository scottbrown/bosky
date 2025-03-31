package beacon

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config holds the application configuration
type Config struct {
	Project    string `yaml:"project"`
	InstanceID string `yaml:"instance_id"`
}

// LoadConfig attempts to load configuration from the given path
// If path is empty, it will check the default location
func LoadConfig(path string) (*Config, error) {
	// Default config
	config := &Config{
		Project: "unknown",
	}

	// If no path specified, try the default location
	if path == "" {
		path = DEFAULT_CONFIG_PATH
	}

	// Check if file exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		data, err := os.ReadFile(path) // #nosec G304
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, err
		}
	} else if path != DEFAULT_CONFIG_PATH {
		// Only return an error if the specified path doesn't exist
		// (but don't error if we just tried the default path)
		return nil, err
	}

	return config, nil
}
