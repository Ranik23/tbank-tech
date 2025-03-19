//go:build unit

package utils

import (
	"encoding/json"
	"testing"

	"github.com/Ranik23/tbank-tech/bot/internal/models"
	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/require"
)

func ConvertFromBytesToCustomCommit_Test(t *testing.T) {

	exampleCustomCommit := models.CustomCommit{
		UserID: 1,
		Commit: &github.RepositoryCommit{
			SHA: github.Ptr("test_sha"),
			Commit: &github.Commit{
				Message: github.Ptr("test_message"),
			},
		},
	}

	bytes, err := json.Marshal(exampleCustomCommit)
	require.NoError(t, err)

	customCommit, err := ConvertFromBytesToCustomCommit(bytes)
	require.NoError(t, err)

	require.Equal(t, customCommit.UserID, exampleCustomCommit.UserID)
	require.Equal(t, customCommit.Commit.SHA, exampleCustomCommit.Commit.SHA)
	require.Equal(t, customCommit.Commit.Commit.Message, exampleCustomCommit.Commit.Commit.Message)
}
