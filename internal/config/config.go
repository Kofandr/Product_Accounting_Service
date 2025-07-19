package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"log"
)

type Configuration struct {
	Port        string `yaml:"Port" env:"PORT" envdefault:"8080"`
	LoggerLevel string `yaml:"LoggerLevel" env:"LOGGER_LEVEL" envdefault:"INFO"`
}

func Load() *Configuration {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading Config.env file: %s", err)
	}
	cfg := &Configuration{}

	if err = env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %s", err)
	}
	return cfg

}
