package utils

import (
	"encoding/json"

	"github.com/Ranik23/tbank-tech/bot/internal/models"
)

func ConvertFromBytesToCustomCommit(message []byte) (*models.CustomCommit, error) {
	var msg models.CustomCommit
	if err := json.Unmarshal(message, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
