package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DB     Database   `envPrefix:"DB_"`
	Server HTTPServer `envPrefix:"SERVER_"`
}

type Database struct {
	Port     int    `env:"PORT"`
	Host     string `env:"HOST"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
	SSLMode  string `env:"SSLMODE"`
}

type HTTPServer struct {
	Addr string `env:"ADDR" envDefault:":8080"`
}

func Load() (*Config, error) {
	var cfg Config

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
