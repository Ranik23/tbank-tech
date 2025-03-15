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



func TestAddLink_Success(t *testing.T) {
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
	repoMock.EXPECT().CreateLink(gomock.Any(), exampleLink.Url).Times(1).Return(nil)
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), exampleLink.Url).Times(1).Return(&exampleLink, nil)
	repoMock.EXPECT().CreateLinkUser(gomock.Any(), exampleLink.ID, exampleID).Times(1).Return(nil)

	hubMock := hub.NewMockHub(ctrl)
	hubMock.EXPECT().AddLink(exampleLink.Url, exampleID).Times(1) 

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	link, err := usecase.AddLink(ctx, exampleLink, exampleID)
	require.NoError(t, err, "Failed to add the link")
	require.Equal(t, exampleLink.Url, link.Url, "Link URL should match")
	require.Equal(t, exampleLink.ID, link.ID, "Link ID should match")
}


func TestAddLink_Fail(t *testing.T) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	_, err = usecase.AddLink(ctx, exampleLink, exampleID)
	require.ErrorIs(t, err, ErrEmptyLink)
}


func TestAddLink_GetTheLinkByURL_Fail(t *testing.T) {
	exampleLink := models.Link{
		Url: "https://github.com/Ranik23/tbank-tech",
	}
	exampleID := uint(1)

	ErrInternal := errors.New("internal error")

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().CreateLink(gomock.Any(), exampleLink.Url).Return(nil).Times(1)
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), exampleLink.Url).Return(nil, ErrInternal).Times(1)

	hubMock := hub.NewMockHub(ctrl)
	hubMock.EXPECT().AddLink(exampleLink.Url, exampleID).Times(1) 

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	_, err = usecase.AddLink(ctx, exampleLink, exampleID)

	require.ErrorIs(t, err, ErrInternal)
}


func TestAddLink_CreateLink_Fail(t *testing.T) {
	exampleLink := models.Link{
		Url: "https://github.com/Ranik23/tbank-tech",
	}
	exampleID := uint(1)

	ErrCreateLink := errors.New("failed to create link")

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().CreateLink(gomock.Any(), exampleLink.Url).Return(ErrCreateLink).Times(1)

	hubMock := hub.NewMockHub(ctrl)
	hubMock.EXPECT().AddLink(exampleLink.Url, exampleID).Times(1) 


	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = usecase.AddLink(ctx, exampleLink, exampleID)
	require.ErrorIs(t, err, ErrCreateLink, "Expected error when creating link")
}


func TestAddLink_CreateLinkUser_Fail(t *testing.T) {
	exampleLink := models.Link{
		Url: "https://github.com/Ranik23/tbank-tech",
	}
	exampleID := uint(1)

	ErrCreateLinkUser := errors.New("failed to create link user")

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().CreateLink(gomock.Any(), exampleLink.Url).Return(nil).Times(1)
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), exampleLink.Url).Return(&exampleLink, nil).Times(1)
	repoMock.EXPECT().CreateLinkUser(gomock.Any(), exampleLink.ID, exampleID).Return(ErrCreateLinkUser).Times(1)

	hubMock := hub.NewMockHub(ctrl)
	hubMock.EXPECT().AddLink(exampleLink.Url, exampleID).Times(1)

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = usecase.AddLink(ctx, exampleLink, exampleID)
	require.ErrorIs(t, err, ErrCreateLinkUser, "Expected error when creating link user")
}

func TestAddLink_Commit_Fail(t *testing.T) {
	exampleLink := models.Link{
		Url: "https://github.com/Ranik23/tbank-tech",
	}
	exampleID := uint(1)

	ErrCommit := errors.New("failed to commit transaction")

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).Times(1)
	mockTx.EXPECT().Commit(gomock.Any()).Return(ErrCommit).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().CreateLink(gomock.Any(), exampleLink.Url).Return(nil).Times(1)
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), exampleLink.Url).Return(&exampleLink, nil).Times(1)
	repoMock.EXPECT().CreateLinkUser(gomock.Any(), exampleLink.ID, exampleID).Return(nil).Times(1)

	hubMock := hub.NewMockHub(ctrl)
	hubMock.EXPECT().AddLink(exampleLink.Url, exampleID).Times(1)

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = usecase.AddLink(ctx, exampleLink, exampleID)
	require.ErrorIs(t, err, ErrCommit, "Expected error when committing transaction")
}


func TestAddLink_Rollback_Fail(t *testing.T) {
	exampleLink := models.Link{
		Url: "https://github.com/Ranik23/tbank-tech",
	}
	exampleID := uint(1)

	ErrRollback := errors.New("rollback failed")
	ErrCommit := errors.New("commit failed")

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Commit(gomock.Any()).Return(ErrCommit).Times(1)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(ErrRollback).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().CreateLink(gomock.Any(), exampleLink.Url).Return(nil).Times(1)
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), exampleLink.Url).Return(&exampleLink, nil).Times(1)
	repoMock.EXPECT().CreateLinkUser(gomock.Any(), exampleLink.ID, exampleID).Return(nil).Times(1)

	hubMock := hub.NewMockHub(ctrl)
	hubMock.EXPECT().AddLink(exampleLink.Url, exampleID).Times(1)

	usecase, err := NewService(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = usecase.AddLink(ctx, exampleLink, exampleID)
	require.ErrorAs(t, err, &ErrRollback, "Expected error when rollback fails")
}

