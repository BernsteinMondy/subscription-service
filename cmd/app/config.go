package main

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		HTTPServer HTTPServer `envPrefix:"HTTP_SERVER_"`
		DB         DB         `envPrefix:"MAIN_DB_"`
		Migrations Migrations `envPrefix:"MIGRATIONS_"`
	}
	HTTPServer struct {
		ListenAddr string `env:"LISTEN_ADDR,notEmpty"`
	}
	DB struct {
		Host         string `env:"HOST,notEmpty"`
		Port         int    `env:"PORT,notEmpty"`
		User         string `env:"USER,notEmpty"`
		Password     string `env:"PASSWORD,notEmpty"`
		DatabaseName string `env:"DATABASE_NAME,notEmpty"`
		SSLMode      string `env:"SSL_MODE,notEmpty"`
	}
	Migrations struct {
		Dir     string `env:"DIR,notEmpty"`
		Enabled bool   `env:"ENABLED" envDefault:"false"`
	}
)

func loadConfigFromEnv() (Config, error) {
	c, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", err)
	}

	return c, nil
}
