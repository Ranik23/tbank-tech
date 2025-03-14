package usecase

import (
	"context"
	"log/slog"
	"tbank/scrapper/internal/mocks/hub"
	//"tbank/scrapper/internal/mocks/pgx/txmanager"
	"tbank/scrapper/internal/mocks/pgx/txunit"
	"tbank/scrapper/internal/mocks/repository"
	"tbank/scrapper/internal/models"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)



func Test(t *testing.T) {
	linkExample := models.Link{
		Url: "https://github.com/Ranik23/tbank-tech",
	}
	userIDExample := 1

	ctrl := gomock.NewController(t)

	mockTx := txunit.NewMockTx(ctrl)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil).Times(1)

	logger := slog.Default()

	repoMock := repository.NewMockRepository(ctrl)
	repoMock.EXPECT().BeginTx(gomock.Any()).Times(1).Return(mockTx, nil)
	repoMock.EXPECT().CreateLink(gomock.Any(), linkExample.Url).Times(1).Return(nil)
	repoMock.EXPECT().GetLinkByURL(gomock.Any(), linkExample.Url).Times(1).Return(&linkExample, nil)
	repoMock.EXPECT().CreateLinkUser(gomock.Any(), linkExample.ID, uint(userIDExample)).Times(1).Return(nil)

	hubMock := hub.NewMockHub(ctrl)
	hubMock.EXPECT().AddLink(linkExample.Url, uint(userIDExample)).Times(1) 

	usecase, err := NewUseCaseImpl(repoMock, hubMock, logger)
	require.NoError(t, err, "Failed to create the usecase")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	link, err := usecase.AddLink(ctx, linkExample, uint(userIDExample))
	require.NoError(t, err, "Failed to add the link")
	require.Equal(t, link.Url, linkExample.Url)
}