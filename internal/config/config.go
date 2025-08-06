package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v8"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Port        int    `env:"PORT" envdefault:"8080" validate:"required,min=1,max=65535"`
	LoggerLevel string `env:"LOGGER_LEVEL" envdefault:"INFO" validate:"required,oneof=DEBUG INFO WARN ERROR"`
	DatabaseUrl string `env:"DATABASE_URL" validate:"required"`
}

func Load() (*Configuration, error) {
	cfg := &Configuration{}

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading Config.env file: %s", err)
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("Failed to parse environment variables: %s", err)
	}

	if err := validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("validator error: %s", err)
	}

	return cfg, nil

}

func Mustload() *Configuration {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}
