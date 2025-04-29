package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port                       int             `yaml:"port"`
	CertFile                   string          `yaml:"cert_file"`
	KeyFile                    string          `yaml:"key_file"`
	TLSEnabled                 bool            `yaml:"tls_enabled"`
	SelectionAlgorithm         string          `yaml:"selection_algorithm"`
	Backends                   []BackendConfig `yaml:"backends"`
	HealthCheckIntervalSeconds int             `yaml:"health_check_interval_seconds"`
}

type BackendConfig struct {
	URL string `yaml:"url"`
}

func LoadConfig(path string) *Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return &cfg
}
