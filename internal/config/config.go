package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MongoDB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
	JWTKey      string `yaml:"jwt_key"`
	Environment string `yaml:"environment"`
	SentryDSN   string `yaml:"sentry_dsn"`
}

func LoadConfig() *Config {
	var configuration Config
	configFile := os.Getenv("FRIEND_MANAGEMENT_CONFIG_FILE")
	if configFile == "" {
		configFile = "config.yaml"
	}

	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(f, &configuration)
	if err != nil {
		panic(err)
	}
	return &configuration
}
