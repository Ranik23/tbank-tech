//go:build integration

package usecasetostorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/Ranik23/tbank-tech/scrapper/config"
	mockhub "github.com/Ranik23/tbank-tech/scrapper/internal/hub/mock"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"
	"github.com/Ranik23/tbank-tech/scrapper/internal/service"

	"github.com/docker/go-connections/nat"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRollBackCheckAddLink(t *testing.T) {

	logger := slog.Default()

	cfg, err := config.LoadConfig("../../../../.env")
	require.NoError(t, err)

	ctx := context.Background()

	exposedPort := fmt.Sprintf("%s/tcp", cfg.DataBase.Port)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{exposedPort},
		WaitingFor:   wait.ForListeningPort(nat.Port(exposedPort)),
		Env: map[string]string{
			"POSTGRES_USER":     cfg.DataBase.Username,
			"POSTGRES_PASSWORD": cfg.DataBase.Password,
			"POSTGRES_DB":       cfg.DataBase.DBName,
		},
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	testcontainers.CleanupContainer(t, postgresC)
	require.NoError(t, err)

	err = postgresC.Start(ctx)
	require.NoError(t, err)

	hostPort, err := postgresC.MappedPort(ctx, nat.Port(exposedPort))
	require.NoError(t, err)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DataBase.Host, hostPort.Port(), cfg.DataBase.Username, cfg.DataBase.Password, cfg.DataBase.DBName, cfg.DataBase.SSL)



	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	sqlDB, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer func() {
		err = sqlDB.Close()
		require.NoError(t, err)
	}()
	err = goose.Up(sqlDB, "../../../internal/migrations")
	require.NoError(t, err)

	txManager := postgres.NewTxManager(pool, logger)
	repository := postgres.NewPostgresRepository(txManager, logger)
	ctrl := gomock.NewController(t)
	mockHub := mockhub.NewMockHub(ctrl)
	mockHub.EXPECT().AddLink(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(errors.New("hub_error"))


	service, err := service.NewService(repository, txManager, mockHub, logger)
	require.NoError(t, err)

	exampleLink := "https://github.com/epchamp001/avito-tech-merch"

	exampleName := "anton"
	exampleID := 1
	exampleToken := "test"

	err = service.RegisterUser(context.Background(), uint(exampleID), exampleName, exampleToken)
	require.NoError(t, err)

	err = service.AddLink(context.Background(), exampleLink, uint(exampleID))
	require.Error(t, err)

	var (
		userID int
		name   string
		link   string
	)

	err = pool.QueryRow(ctx, `SELECT user_id, name FROM users WHERE user_id = $1`, exampleID).Scan(&userID, &name)
	require.NoError(t, err)
	require.Equal(t, exampleID, userID)
	require.Equal(t, exampleName, name)

	err = pool.QueryRow(ctx, `SELECT url FROM links WHERE url = $1`, exampleLink).Scan(&link)
	require.Error(t, err)
}


func TestRollBackCheckRemoveLink(t *testing.T) {
	logger := slog.Default()

	cfg, err := config.LoadConfig("../../../../.env")
	require.NoError(t, err)

	ctx := context.Background()

	exposedPort := fmt.Sprintf("%s/tcp", cfg.DataBase.Port)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{exposedPort},
		WaitingFor:   wait.ForListeningPort(nat.Port(exposedPort)),
		Env: map[string]string{
			"POSTGRES_USER":     cfg.DataBase.Username,
			"POSTGRES_PASSWORD": cfg.DataBase.Password,
			"POSTGRES_DB":       cfg.DataBase.DBName,
		},
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	testcontainers.CleanupContainer(t, postgresC)
	require.NoError(t, err)

	err = postgresC.Start(ctx)
	require.NoError(t, err)

	hostPort, err := postgresC.MappedPort(ctx, nat.Port(exposedPort))
	require.NoError(t, err)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DataBase.Host, hostPort.Port(), cfg.DataBase.Username, cfg.DataBase.Password, cfg.DataBase.DBName, cfg.DataBase.SSL)



	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	sqlDB, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer func() {
		err = sqlDB.Close()
		require.NoError(t, err)
	}()
	err = goose.Up(sqlDB, "../../../internal/migrations")
	require.NoError(t, err)

	txManager := postgres.NewTxManager(pool, logger)
	repository := postgres.NewPostgresRepository(txManager, logger)
	ctrl := gomock.NewController(t)
	mockHub := mockhub.NewMockHub(ctrl)

	mockHub.EXPECT().RemoveLink(gomock.Any(), gomock.Any()).AnyTimes().Return(errors.New("hub error"))
	
	mockHub.EXPECT().AddLink(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)


	service, err := service.NewService(repository, txManager, mockHub, logger)
	require.NoError(t, err)


	exampleLink := "https://github.com/epchamp001/avito-tech-merch"

	exampleName := "anton"
	exampleID := 1
	exampleToken := "test"

	err = service.RegisterUser(context.Background(), uint(exampleID), exampleName, exampleToken)
	require.NoError(t, err)

	err = service.AddLink(context.Background(), exampleLink, uint(exampleID))
	require.NoError(t, err)
	

	// тут будет вызван роллбек, тк транзакция завершится с ошибкой от mockHub 
	// и defer обязан будет поймать ошибку и сделать роллбек транзакции
	err = service.RemoveLink(context.Background(), exampleLink, uint(exampleID))
	require.Error(t, err)

	var (
		link string
	)

	// обязан найти, потому что мы сделали роллбек и все удаленное было возвращено
	// причем роллбек был вызван после того как запись была удалена
	err = pool.QueryRow(ctx, `SELECT url FROM links WHERE url = $1`, exampleLink).Scan(&link)
	require.NoError(t, err)
}
