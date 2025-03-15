package service


import (
	"context"
	"errors"
	"log/slog"
	"tbank/scrapper/internal/mocks/hub"
	"tbank/scrapper/internal/mocks/pgx/txunit"
	"tbank/scrapper/internal/mocks/repository"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestRegisterUser_Success(t *testing.T) {
	exampleUserID := uint(1)
	exampleName := "John Doe"

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().CreateUser(gomock.Any(), exampleUserID, exampleName).Times(1).Return(nil)

	hubMock := hub.NewMockHub(ctrl)

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = usecase.RegisterUser(ctx, exampleUserID, exampleName)
	require.NoError(t, err, "Failed to register the user")
}

func TestRegisterUser_Fail_CreateUser_Error(t *testing.T) {
	exampleUserID := uint(1)
	exampleName := "John Doe"

	ErrCreateUser := errors.New("create user failed")

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().CreateUser(gomock.Any(), exampleUserID, exampleName).Times(1).Return(ErrCreateUser)

	hubMock := hub.NewMockHub(ctrl)

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверка ошибки при создании пользователя
	err = usecase.RegisterUser(ctx, exampleUserID, exampleName)
	require.ErrorIs(t, err, ErrCreateUser)
}

func TestRegisterUser_Fail_RollbackError(t *testing.T) {
	exampleUserID := uint(1)
	exampleName := "John Doe"

	ErrCreateUser := errors.New("create user failed")
	ErrRollback := errors.New("rollback failed")

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(ErrRollback).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().CreateUser(gomock.Any(), exampleUserID, exampleName).Times(1).Return(ErrCreateUser)

	hubMock := hub.NewMockHub(ctrl)

	usecase, err := NewService(repoMock,hubMock,  logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверка ошибки при откате транзакции
	err = usecase.RegisterUser(ctx, exampleUserID, exampleName)
	require.ErrorAs(t, err, &ErrRollback, "Expected rollback error")
}
