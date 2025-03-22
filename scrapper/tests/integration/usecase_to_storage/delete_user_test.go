//go:build integration

package usecasetostorage

import (
	"context"
	"database/sql"
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

func TestDeletesUser(t *testing.T) {

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

	mockHub.EXPECT().AddLink(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	serv, err := service.NewService(repository, txManager, mockHub, logger)
	require.NoError(t, err)

	err = serv.DeleteUser(ctx, 1)
	require.Error(t, err)

	err = serv.RegisterUser(ctx, 45, "egor", "test")
	require.NoError(t, err)

	err = serv.DeleteUser(ctx, 45)
	require.NoError(t, err)

	_, err = serv.GetLinks(ctx, 45)
	require.Error(t, err)
}
