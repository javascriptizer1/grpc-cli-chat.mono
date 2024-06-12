package config

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `env:"ENV" env-default:"local"`
	GRPC       GRPCConfig
	GRPCClient GRPCClientConfig
	HTTP       HTTPConfig
	DB         DBConfig
}

type GRPCConfig struct {
	Host    string        `env:"GRPC_SERVER_HOST" env-default:"localhost"`
	Port    int           `env:"GRPC_SERVER_PORT" env-default:"50051"`
	Timeout time.Duration `env:"GRPC_SERVER_TIMEOUT"`
}

func (c *GRPCConfig) HostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type GRPCClientConfig struct {
	Host string `env:"GRPC_CLIENT_HOST" env-default:"localhost"`
	Port int    `env:"GRPC_CLIENT_PORT" env-default:"50052"`
}

func (c *GRPCClientConfig) HostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type HTTPConfig struct {
	Host    string        `env:"HTTP_SERVER_HOST" env-default:"localhost"`
	Port    int           `env:"HTTP_SERVER_PORT" env-default:"8081"`
	Timeout time.Duration `env:"HTTP_SERVER_TIMEOUT"`
}

func (c *HTTPConfig) HostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type DBConfig struct {
	Host     string `env:"MONGO_HOST" env-required:"true"`
	Port     string `env:"MONGO_PORT" env-required:"true"`
	User     string `env:"MONGO_USER" env-required:"true"`
	Password string `env:"MONGO_PASSWORD" env-required:"true"`
	Name     string `env:"MONGO_DB" env-required:"true"`
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
