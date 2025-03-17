package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DataBaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSL      string
}

func (ds *DataBaseConfig) ConnectToPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	if ds.Host == "" || ds.Port == "" || ds.Username == "" || ds.Password == "" || ds.DBName == "" {
		return nil, fmt.Errorf("missing required database connection parameters")
	}
	connectionString := ds.getDSN()
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		return nil, err
	}
	return pool, nil
}

func (ds *DataBaseConfig) getDSN() string {
	if ds.SSL == "" {
		ds.SSL = "disable"
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		ds.Host, ds.Port, ds.Username, ds.Password, ds.DBName, ds.SSL)
	return dsn
}
