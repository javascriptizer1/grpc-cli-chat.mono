package config

import (
	"log"
	"net"
	"path/filepath"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	Version  string `env:"VERSION" env-required:"true"`
	GRPCAuth GRPCAuthConfig
	GRPCChat GRPChatConfig
}

func (c *Config) ClientConfigPath() string {
	return filepath.Join(".gchat", ".config.json")
}

type GRPCAuthConfig struct {
	Host string `env:"GRPC_AUTH_HOST" env-default:"localhost"`
	Port int    `env:"GRPC_AUTH_PORT" env-default:"50051"`
}

func (c *GRPCAuthConfig) HostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type GRPChatConfig struct {
	Host string `env:"GRPC_CHAT_HOST" env-default:"localhost"`
	Port int    `env:"GRPC_CHAT_PORT" env-default:"50052"`
}

func (c *GRPChatConfig) HostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
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
