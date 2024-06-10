package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env  string `env:"ENV" env-default:"local"`
	GRPC GRPCConfig
	DB   DBConfig
}

type GRPCConfig struct {
	Port    int           `env:"GRPC_SERVER_PORT" env-default:"50051"`
	Timeout time.Duration `env:"GRPC_SERVER_TIMEOUT"`
}

type DBConfig struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     string `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Name     string `env:"POSTGRES_DB" env-required:"true"`
}

func MustLoad() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("config is empty: " + err.Error())
	}

	return &cfg
}
