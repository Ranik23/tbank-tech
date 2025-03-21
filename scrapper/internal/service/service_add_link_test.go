//go:build unit

package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	hubmock "github.com/Ranik23/tbank-tech/scrapper/internal/hub/mock"
	dbmodels "github.com/Ranik23/tbank-tech/scrapper/internal/models"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/mock"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)



func TestAddLinkSuccess(t *testing.T) {

	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	link := "https://github.com/epchamp001/avito-tech-merch"
	userID := uint(1)
	user := &dbmodels.User{UserID: uint(1), Name: "anton", Token: "user-token"}

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn)).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	repoMock.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), link).Return(nil, postgres.ErrNoLinkFound)

	// Создаем ссылку
	repoMock.EXPECT().CreateLink(gomock.Any(), link).Return(nil)

	// Получаем ссылку после создания
	linkObj := &dbmodels.Link{ID: 1, Url: link}
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), link).Return(linkObj, nil)

	// Создаем связь между пользователем и ссылкой
	repoMock.EXPECT().CreateLinkUser(gomock.Any(), linkObj.ID, userID).Return(nil)

	hubMock.EXPECT().AddLink(link, userID, user.Token, 10 * time.Second).Return(nil)

	err = service.AddLink(context.Background(), link, userID)
	require.NoError(t, err)
}


func TestAddLinkUserNotFound(t *testing.T) {
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
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn)).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	repoMock.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, postgres.ErrNoUserFound)

	err = service.AddLink(context.Background(), link, userID)
	require.ErrorIs(t, err, ErrUserNotFound)
}

func TestAddLinkDbErrorOnCreateLinkUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := mock.NewMockRepository(ctrl)
	txManagerMock := mock.NewMockTxManager(ctrl)
	hubMock := hubmock.NewMockHub(ctrl)
	logger := slog.Default()

	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
	require.NoError(t, err)

	link := "https://github.com/epchamp001/avito-tech-merch"
	userID := uint(1)
	user := &dbmodels.User{UserID: userID, Token: "user-token"}

	var fn func(ctx context.Context) error
	txManagerMock.EXPECT().WithTx(gomock.Any(), gomock.AssignableToTypeOf(fn)).
	DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})

	repoMock.EXPECT().GetUserByID(gomock.Any(), userID).Return(user, nil)
	
	linkObj := &dbmodels.Link{ID: 1, Url: link}
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), link).Return(linkObj, nil)

	// Ошибка при создании связи
	repoMock.EXPECT().CreateLinkUser(gomock.Any(), linkObj.ID, userID).Return(errors.New("failed to link user and link"))

	err = service.AddLink(context.Background(), link, userID)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to link user and link")
}

