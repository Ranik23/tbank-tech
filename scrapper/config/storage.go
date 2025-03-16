package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)


type DataBaseConfig struct {		
	Host 		string				
	Port 		string			
	Username 	string				
	Password 	string					
	DBName 		string
	SSL			string					
}

func (ds *DataBaseConfig) ConnectToPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	connectionString := ds.getDSN()
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (ds *DataBaseConfig) getDSN() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		ds.Host, ds.Port, ds.Username, ds.Password, ds.DBName, ds.SSL)
	return dsn
}
