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

func TestGetLinks(t *testing.T) {
	// Настройка логгера
	logger := slog.Default()

	cfg, err := config.LoadConfig("/home/anton/tbank-tech/.env")
	require.NoError(t, err)

	// Настройка контейнера с PostgreSQL
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
	require.NoError(t, err)
	defer func() {
		testcontainers.CleanupContainer(t, postgresC)
	}()

	err = postgresC.Start(ctx)
	require.NoError(t, err)

	hostPort, err := postgresC.MappedPort(ctx, nat.Port(exposedPort))
	require.NoError(t, err)

	// Подключение к базе данных
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

	// Применение миграций
	err = goose.Up(sqlDB, "../../../internal/migrations")
	require.NoError(t, err)

	// Создание сервиса
	txManager := postgres.NewTxManager(pool, logger)
	repository := postgres.NewPostgresRepository(txManager, logger)

	ctrl := gomock.NewController(t)
	mockHub := mockhub.NewMockHub(ctrl)
	mockHub.EXPECT().AddLink(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	serv, err := service.NewService(repository, txManager, mockHub, logger)
	require.NoError(t, err)

	// Подготовка данных для тестов
	_, err = pool.Exec(ctx, `INSERT INTO users (user_id, name, token) VALUES ($1, $2, $3)`, 1, "test_user", "test_token")
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `INSERT INTO links (url) VALUES ($1)`, "https://example.com")
	require.NoError(t, err)

	var linkID uint
	err = pool.QueryRow(ctx, `SELECT id FROM links WHERE url = $1`, "https://example.com").Scan(&linkID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `INSERT INTO link_users (link_id, user_id) VALUES ($1, $2)`, linkID, 1)
	require.NoError(t, err)

	// Таблица тестовых случаев
	testCases := []struct {
		name          string
		userID        uint
		expectedLinks int
		expectedError error
	}{
		{
			name:          "Non-existent user2",
			userID:        2, // Пользователь без ссылок
			expectedLinks: 0,
			expectedError: service.ErrUserNotFound,
		},
		{
			name:          "User with links",
			userID:        1, // Пользователь с ссылками
			expectedLinks: 1,
			expectedError: nil,
		},
		{
			name:          "Non-existent user",
			userID:        999, // Несуществующий пользователь
			expectedLinks: 0,
			expectedError: service.ErrUserNotFound,
		},
	}

	// Запуск тестовых случаев
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			links, err := serv.GetLinks(ctx, tc.userID)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedLinks, len(links))
			}
		})
	}
}
