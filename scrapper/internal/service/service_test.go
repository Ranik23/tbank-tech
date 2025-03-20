package service

// import (
// 	"context"
// 	"log/slog"
// 	"testing"

// 	hubmock "github.com/Ranik23/tbank-tech/scrapper/internal/hub/mock"
// 	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/mock"
// 	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/require"
// )




// func TestRegisterUser(t *testing.T) {

// 	ctrl := gomock.NewController(t)
// 	repoMock := mock.NewMockRepository(ctrl)
// 	txManagerMock := postgres.NewTxManager()
// 	hubMock := hubmock.NewMockHub(ctrl)
// 	logger := slog.Default()

// 	service, err := NewService(repoMock, txManagerMock, hubMock, logger)
// 	require.NoError(t, err)

// 	repoMock.EXPECT().GetUserByID(gomock.Any(), gomock.Eq(1)).Return(nil, postgres.ErrNoUserFound).Times(1)
// 	repoMock.EXPECT().CreateUser(gomock.Any(), gomock.Eq(1), "anton", "token").Return(nil).Times(1)

// 	err = service.RegisterUser(context.Background(), 1, "anton", "token")
// 	require.NoError(t, err)
// }