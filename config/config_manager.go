package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// MongoConfig has config values for MongoDB
type MongoConfig string

// LogentriesConfig has config values for Logentries
type LogentriesConfig string

// RedisConfig has config values for Redis
type RedisConfig struct {
	Addr                string `json:"addr"`
	Pwd                 string `json:"pwd"`
	PendingAlertsKey    string `json:"pendingAlertsKey"`
	ProcessingAlertsKey string `json:"processingAlertsKey"`
}

// Config struct defines the config structure
type Config struct {
	Mongo      MongoConfig      `json:"mongo"`
	Redis      RedisConfig      `json:"redis"`
	Logentries LogentriesConfig `json:"logentries"`
	Secret     string           `json:"secret"`
}

// NewConfig parses config file and return Config struct
func NewConfig(configPath string) *Config {
	var file []byte
	var err error

	if configPath != "" {
		file, err = ioutil.ReadFile(configPath)

		if err != nil {
			log.Fatalf("Config file '%s' file not found", configPath)
		}
	} else {
		file, err = ioutil.ReadFile("./config.json")

		if err != nil {
			file, err = ioutil.ReadFile("/etc/finch/config.json")

			if err != nil {
				log.Fatalf("Config file is not found")
			}
		}
	}

	config := &Config{}
	err = json.Unmarshal(file, config)

	if err != nil {
		panic(err)
	}

	return config
}
