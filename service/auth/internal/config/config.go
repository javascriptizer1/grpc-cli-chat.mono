package config

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env  string `env:"ENV" env-default:"local"`
	GRPC GRPCConfig
	HTTP HTTPConfig
	DB   DBConfig
	JWT  JWTConfig
}

type GRPCConfig struct {
	Host    string        `env:"GRPC_SERVER_HOST" env-default:"localhost"`
	Port    int           `env:"GRPC_SERVER_PORT" env-default:"50051"`
	Timeout time.Duration `env:"GRPC_SERVER_TIMEOUT"`
}

func (c *GRPCConfig) HostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type HTTPConfig struct {
	Host    string        `env:"HTTP_SERVER_HOST" env-default:"localhost"`
	Port    int           `env:"HTTP_SERVER_PORT" env-default:"8080"`
	Timeout time.Duration `env:"HTTP_SERVER_TIMEOUT"`
}

func (c *HTTPConfig) HostPort() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
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
	var cfg Config
	var err error

	configPath := fetchConfigPath()

	if configPath != "" {
		err = godotenv.Load(configPath)
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Printf("No loading .env file: %v", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("config is empty: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
