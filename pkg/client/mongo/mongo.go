package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (c *Config) toURL() string {
	if c.Port == "" {
		return fmt.Sprintf("mongodb+srv://%s:%s@%s/?tls=true&tlsInsecure=true", c.User, c.Password, c.Host)
	}

	return fmt.Sprintf("mongodb://%s:%s@%s:%s", c.User, c.Password, c.Host, c.Port)
}

func New(ctx context.Context, config Config) (*mongo.Database, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(config.toURL()).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}

	db := client.Database(config.Name)

	return db, nil
}
