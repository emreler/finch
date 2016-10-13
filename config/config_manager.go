package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const configFile = "config.json"

// MongoConfig has config values for MongoDB
type MongoConfig string

// RedisConfig has config values for Redis
type RedisConfig struct {
	Addr string `json:"addr"`
	Pwd  string `json:"pwd"`
}

// Config struct defines the config structure
type Config struct {
	Mongo MongoConfig `json:"mongo"`
	Redis RedisConfig `json:"redis"`
}

// NewConfig parses config file and return Config struct
func NewConfig() *Config {
	file, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatalf("Config file '%s' file not found", configFile)
	}

	config := &Config{}
	json.Unmarshal(file, config)

	return config
}
