package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (c *Config) toURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.Name)
}

func New(config Config) (*sqlx.DB, error) {
	connConfig, _ := pgx.ParseConfig(config.toURL())
	connStr := stdlib.RegisterConnConfig(connConfig)

	db, err := sqlx.Open("pgx", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
