package fixss

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	Logging struct {
		Config string `yaml:"config"`
	}
	Fix struct {
		Config string `yaml:"config"`
	}
	Quote struct {
		Config string `yaml:"defaultConfig"`
	}
}

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
