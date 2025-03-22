//go:build unit

package service

import (
	"context"
	"errors"

	//"errors"
	"log/slog"
	"testing"

	hubmock "github.com/Ranik23/tbank-tech/scrapper/internal/hub/mock"
	"github.com/jackc/pgx/v5"
	//"github.com/Ranik23/tbank-tech/scrapper/internal/models"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/mock"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"

	//"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestDeleteUserSucces(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn), gomock.Any()).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error, mode pgx.TxAccessMode) error {
		return fn(ctx)
	})

	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, nil)
	repoMock.EXPECT().DeleteUser(gomock.Any(), uint(1)).Return(nil)


	err = service.DeleteUser(context.Background(), 1)
	require.NoError(t, err)
}

func TestDeleteUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn), gomock.Any()).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error, mode pgx.TxAccessMode) error {
		return fn(ctx)
	})

	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, postgres.ErrNoUserFound)

	err = service.DeleteUser(context.Background(), 1)
	require.ErrorIs(t, err, ErrUserNotFound)
}


func TestDeleteUserDbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn), gomock.Any()).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error, mode pgx.TxAccessMode) error {
		return fn(ctx)
	})

	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, errors.New("db connection error"))

	err = service.DeleteUser(context.Background(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db connection error")
}

func TestDeleteUserDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn), gomock.Any()).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error, mode pgx.TxAccessMode) error {
		return fn(ctx)
	})

	// Симулируем, что пользователь найден
	repoMock.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(nil, nil)
	// Симулируем, что удаление не удалось
	repoMock.EXPECT().DeleteUser(gomock.Any(), uint(1)).Return(errors.New("failed to delete user"))

	err = service.DeleteUser(context.Background(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to delete user")
}

func TestDeleteUserTxError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(func(ctx context.Context) error { return nil }), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error, mode pgx.TxAccessMode) error {
			return errors.New("transaction error")
		})

	err = service.DeleteUser(context.Background(), 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction error")
}

