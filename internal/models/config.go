package models

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password,omitempty"`
	Key      string `yaml:"key,omitempty"`
}

type Config struct {
	Servers []Server `yaml:"servers"`
}

func LoadConfig(filePath string) (*Config, error) {
	if len(filePath) == 0 {
		filePath = os.Getenv("SENTRY_CONFIG_PATH")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
