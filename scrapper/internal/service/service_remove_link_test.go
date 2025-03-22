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
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)


func TestRemoveLinkSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	link := "https://github.com/epchamp001/avito-tech-merch"
	userID := uint(1)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn), gomock.Any()).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error, mode pgx.TxAccessMode) error {
		return fn(ctx)
	})

	// Проверяем, что ссылка есть
	linkObj := &dbmodels.Link{ID: 1, Url: link}
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), link).Return(linkObj, nil)

	// Проверяем, что пользователь существует
	repoMock.EXPECT().GetUserByID(gomock.Any(), userID).Return(&dbmodels.User{UserID: userID}, nil)

	// Удаляем ссылку
	repoMock.EXPECT().DeleteLink(gomock.Any(), linkObj.ID).Return(nil)

	// Убираем из Hub
	hubMock.EXPECT().RemoveLink(link, userID).Return(nil)

	err = service.RemoveLink(context.Background(), link, userID)
	require.NoError(t, err)
}


func TestRemoveLinkNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	
	link := "https://github.com/epchamp001/avito-tech-merch"
	userID := uint(1)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn), gomock.Any()).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error, mode pgx.TxAccessMode) error {
		return fn(ctx)
	})

	// Симулируем, что ссылка не найдена
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), link).Return(nil, postgres.ErrNoLinkFound)

	err = service.RemoveLink(context.Background(), link, userID)
	require.Error(t, err)
	require.Equal(t, ErrLinkNotFound, err)
}


func TestRemoveLinkUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	link := "https://github.com/epchamp001/avito-tech-merch"
	userID := uint(1)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn), gomock.Any()).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error, mode pgx.TxAccessMode) error {
		return fn(ctx)
	})

	// Проверяем, что ссылка есть
	linkObj := &dbmodels.Link{ID: 1, Url: link}
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), link).Return(linkObj, nil)

	// Симулируем, что пользователь не найден
	repoMock.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, postgres.ErrNoUserFound)

	err = service.RemoveLink(context.Background(), link, userID)
	require.Error(t, err)
	require.Equal(t, ErrUserNotFound, err)
}

func TestRemoveLinkDeleteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	link := "https://github.com/epchamp001/avito-tech-merch"
	userID := uint(1)

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn), gomock.Any()).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error, mode pgx.TxAccessMode) error {
		return fn(ctx)
	})

	// Проверяем, что ссылка есть
	linkObj := &dbmodels.Link{ID: 1, Url: link}
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), link).Return(linkObj, nil)

	// Проверяем, что пользователь существует
	repoMock.EXPECT().GetUserByID(gomock.Any(), userID).Return(&dbmodels.User{UserID: userID}, nil)

	// Симулируем ошибку при удалении
	repoMock.EXPECT().DeleteLink(gomock.Any(), linkObj.ID).Return(errors.New("delete error"))

	err = service.RemoveLink(context.Background(), link, userID)
	require.Error(t, err)
}


