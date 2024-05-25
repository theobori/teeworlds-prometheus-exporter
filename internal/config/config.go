package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Servers Servers `yaml:"servers"`
}

type Servers struct {
	Econ   []EconServer   `yaml:"econ"`
	Master []MasterServer `yaml:"master"`
}

type EconServer struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Password string `yaml:"password"`
}

type MasterServer struct {
	Protocol        string `yaml:"protocol"`
	URL             string `yaml:"url,omitempty"`
	Host            string `yaml:"host,omitempty"`
	Port            uint16 `yaml:"port,omitempty"`
	RefreshCooldown uint   `yaml:"refresh_cooldown" default:"10"`
}

// Get YAML data as `Config`
func ConfigFromData(data []byte) (*Config, error) {
	var config Config

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Get YAML data as `Config` from a file
func ConfigFromFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return ConfigFromData(data)
}
