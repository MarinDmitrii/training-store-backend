package builder

import "os"

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	password string
	Database string
}

func NewPostgresConfig() (PostgresConfig, error) {
	return PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DATABASE"),
	}, nil
}
