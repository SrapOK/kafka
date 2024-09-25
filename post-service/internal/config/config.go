package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServer `yaml:"http_server"`
	Kafka      `yaml:"kafka"`
}

type HttpServer struct {
	Addr        string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"8s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" evn-default:"10s"`
}

type Kafka struct {
	Addr  string `yaml:"address" env-default:"localhost:9092"`
	Topic string `yaml:"topic"`
}

func Load(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	return &cfg
}
