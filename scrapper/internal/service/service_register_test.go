//go:build unit

package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	hubmock "github.com/Ranik23/tbank-tech/scrapper/internal/hub/mock"
	"github.com/Ranik23/tbank-tech/scrapper/internal/models"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/mock"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestRegisterUserSucces(t *testing.T) {

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

	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, postgres.ErrNoUserFound).Times(1)
	repoMock.EXPECT().CreateUser(gomock.Any(), uint(1), "anton", "token").Return(nil).Times(1)

	err = service.RegisterUser(context.Background(), 1, "anton", "token")
	require.NoError(t, err)
}


func TestRegisterUser_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).
		Return(&models.User{UserID: 1, Name: "anton", Token: "token"}, nil).Times(1) // Пользователь найден

	err = service.RegisterUser(context.Background(), 1, "anton", "token")
	require.ErrorIs(t, err, ErrUserAlreadyExists)
}

func TestRegisterUser_CreateUserFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)

	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).
		Return(nil, postgres.ErrNoUserFound).Times(1) // Пользователь не найден

	repoMock.EXPECT().CreateUser(gomock.Any(), uint(1), "anton", "token").
		Return(errors.New("db error")).Times(1) // Ошибка БД

	err = service.RegisterUser(context.Background(), 1, "anton", "token")
	require.ErrorContains(t, err, "db error") // Проверяем текст ошибки
}


func TestRegisterUser_WithTxFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.Any()).Return(errors.New("tx failed"))

	err = service.RegisterUser(context.Background(), 1, "anton", "token")
	require.ErrorContains(t, err, "tx failed")
}
