//go:build integration

package usecasestorage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"tbank/scrapper/config"
	mockhub "tbank/scrapper/internal/hub/mock"
	"tbank/scrapper/internal/models"
	"tbank/scrapper/internal/repository/postgres"
	"tbank/scrapper/internal/service"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Test(t *testing.T) {

	logger := slog.Default()

	_, currentFile, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(currentFile)

	cfg, err := config.LoadConfig(filepath.Join(testDir, ".env"))

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

	err = goose.Up(sqlDB, "../../../internal/migrations")
	require.NoError(t, err)


	txManager := postgres.NewTxManager(pool, logger)

	repository := postgres.NewPostgresRepository(txManager, logger)

	ctrl := gomock.NewController(t)

	mockHub := mockhub.NewMockHub(ctrl)

	mockHub.EXPECT().AddLink(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	serv, err := service.NewService(repository, txManager, mockHub, logger)
	require.NoError(t, err)

	exampleLink := models.Link{
		ID:  1,
		Url: "https://github.com/epchamp001/avito-tech-merch",
	}

	_, err = serv.AddLink(context.Background(), exampleLink, 1)
	require.NoError(t, err)

	var userID int
	var name string

	err = pool.QueryRow(ctx, `SELECT user_id, name FROM users WHERE user_id = $1`, 1).Scan(&userID, &name)
	require.NoError(t, err)

	require.Equal(t, 1, userID)
	require.Equal(t, "random", name)

	var link string

	err = pool.QueryRow(ctx, `SELECT url FROM links WHERE url = $1`, "https://github.com/epchamp001/avito-tech-merch").Scan(&link)
	require.NoError(t, err)

	require.Equal(t, link, "https://github.com/epchamp001/avito-tech-merch")
}
