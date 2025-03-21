//go:build unit

package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	hubmock "github.com/Ranik23/tbank-tech/scrapper/internal/hub/mock"
	dbmodels "github.com/Ranik23/tbank-tech/scrapper/internal/models"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/mock"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)




func TestGetLinksSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn)).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	// Пользователь существует
	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, nil)

	// у пользователя будут ссылки
	expectedLinks := []dbmodels.Link{
		{ID: 1, Url: "https://example.com"},
		{ID: 2, Url: "https://golang.org"},
	}
	repoMock.EXPECT().GetLinks(gomock.Any(), uint(1)).Return(expectedLinks, nil)

	links, err := service.GetLinks(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, expectedLinks, links)
}

func TestGetLinksUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn)).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, postgres.ErrNoUserFound)

	links, err := service.GetLinks(context.Background(), 1)
	require.Nil(t, links)
	require.ErrorIs(t, err, ErrUserNotFound)
}


func TestGetLinksDbErrorOnGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn)).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	// Ошибка базы данных при получении пользователя
	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, errors.New("database is down"))

	links, err := service.GetLinks(context.Background(), 1)
	require.Nil(t, links)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get user")
}

func TestGetLinksDbErrorOnGetLinks(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn)).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	// Пользователь найден
	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&dbmodels.User{UserID: 1}, nil)

	// Ошибка при получении ссылок
	repoMock.EXPECT().GetLinks(gomock.Any(), uint(1)).Return(nil, errors.New("failed to fetch links"))

	links, err := service.GetLinks(context.Background(), 1)
	require.Nil(t, links)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get links")
}



func TestGetLinksTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	// Ошибка при начале транзакции
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(func(ctx context.Context) error { return nil })).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return errors.New("transaction error")
		})

	links, err := service.GetLinks(context.Background(), 1)
	require.Nil(t, links)
	require.Error(t, err)
	require.Equal(t, err.Error(), "transaction error")
}


