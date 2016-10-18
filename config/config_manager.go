package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const defaultConfigPath = "config.json"

// MongoConfig has config values for MongoDB
type MongoConfig string

// LogentriesConfig has config values for Logentries
type LogentriesConfig string

// RedisConfig has config values for Redis
type RedisConfig struct {
	Addr string `json:"addr"`
	Pwd  string `json:"pwd"`
}

// Config struct defines the config structure
type Config struct {
	Mongo      MongoConfig      `json:"mongo"`
	Redis      RedisConfig      `json:"redis"`
	Logentries LogentriesConfig `json:"Logentries"`
	Secret     string           `json:"secret"`
}

// NewConfig parses config file and return Config struct
func NewConfig(configPath string) *Config {
	if configPath == "" {
		configPath = defaultConfigPath
	}

	file, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatalf("Config file '%s' file not found", configPath)
	}

	config := &Config{}
	json.Unmarshal(file, config)

	return config
}
