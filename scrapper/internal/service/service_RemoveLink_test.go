package service

import (
	"context"
	"errors"
	"log/slog"
	"tbank/scrapper/internal/mocks/hub"
	"tbank/scrapper/internal/mocks/pgx/txunit"
	"tbank/scrapper/internal/mocks/repository"
	"tbank/scrapper/internal/models"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRemoveLink_Success(t *testing.T) {
	exampleLink := models.Link{
		Url: "https://github.com/Ranik23/tbank-tech",
	}
	exampleID := uint(1)

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), exampleLink.Url).Times(1).Return(&exampleLink, nil)
	repoMock.EXPECT().DeleteLink(gomock.Any(), exampleLink.ID).Times(1).Return(nil)

	hubMock := hub.NewMockHub(ctrl)
	hubMock.EXPECT().RemoveLink(exampleLink.Url, exampleID).Times(1)

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = usecase.RemoveLink(ctx, exampleLink, exampleID)
	require.NoError(t, err, "Failed to remove the link")
}

func TestRemoveLink_Fail_EmptyLinkURL(t *testing.T) {
	exampleLink := models.Link{
		Url: "",
	}
	exampleID := uint(1)

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)

	hubMock := hub.NewMockHub(ctrl)

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = usecase.RemoveLink(ctx, exampleLink, exampleID)
	require.ErrorIs(t, err, ErrEmptyLink)
}

func TestRemoveLink_Fail_RollbackError(t *testing.T) {
	exampleLink := models.Link{
		Url: "https://github.com/Ranik23/tbank-tech",
	}
	exampleID := uint(1)

	ErrInternal := errors.New("internal error")
	ErrRollback := errors.New("rollback failed")

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(ErrRollback).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), exampleLink.Url).Times(1).Return(&exampleLink, nil)
	repoMock.EXPECT().DeleteLink(gomock.Any(), exampleLink.ID).Times(1).Return(ErrInternal)

	hubMock := hub.NewMockHub(ctrl)
	hubMock.EXPECT().RemoveLink(exampleLink.Url, exampleID).Times(1)

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = usecase.RemoveLink(ctx, exampleLink, exampleID)
	require.ErrorAs(t, err, &ErrRollback, "Expected rollback error")
}
