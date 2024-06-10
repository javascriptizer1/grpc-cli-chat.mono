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
	JWT  JWTConfig
}

type GRPCConfig struct {
	Port    int           `env:"GRPC_SERVER_PORT" env-default:"50051"`
	Timeout time.Duration `env:"GRPC_SERVER_TIMEOUT"`
}

type JWTConfig struct {
	AccessSecretKey  string        `env:"JWT_ACCESS_SECRET" env-required:"true"`
	AccessDuration   time.Duration `env:"JWT_ACCESS_DURATION" env-required:"true"`
	RefreshSecretKey string        `env:"JWT_REFRESH_SECRET" env-required:"true"`
	RefreshDuration  time.Duration `env:"JWT_REFRESH_DURATION" env-required:"true"`
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
