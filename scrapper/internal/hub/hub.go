package hub

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	git "tbank/scrapper/pkg/github"
	"time"

	"github.com/IBM/sarama"
)

var (
	ErrEmptyLink = fmt.Errorf("empty link")
)

type Hub struct {
	linksCancel   		map[string]context.CancelFunc
	linksUsers    		map[string][]uint
	kafkaProducer 		sarama.AsyncProducer
	mut           		sync.Mutex
	logger        		*slog.Logger
	gitClient     		git.GitHubClient
}

func NewHub(producer sarama.AsyncProducer, logger *slog.Logger, gitClient git.GitHubClient) *Hub {
	return &Hub{
		linksCancel:   make(map[string]context.CancelFunc),
		linksUsers:    make(map[string][]uint),
		kafkaProducer: producer,
		logger:        logger,
		gitClient:     gitClient,
	}
}

func (h *Hub) AddTrack(linkToTrack string, userID uint) error {
	h.mut.Lock()
	defer h.mut.Unlock()

	if _, exists := h.linksUsers[linkToTrack]; !exists {

		ctx, cancel := context.WithCancel(context.Background())

		client := NewClient(h.kafkaProducer, h.logger, linkToTrack, h.gitClient)

		go client.Run(ctx, 10 * time.Second)

		h.linksCancel[linkToTrack] = cancel
	}

	if !slices.Contains(h.linksUsers[linkToTrack], userID) {
		h.linksUsers[linkToTrack] = append(h.linksUsers[linkToTrack], userID)
	}
	
	return nil
}

func (h *Hub) RemoveTrack(linkToUnTrack string, userID uint) error {
	h.mut.Lock()
	defer h.mut.Unlock()

	if users, exists := h.linksUsers[linkToUnTrack]; exists {

		filteredUsers := []uint{}

		for _, id := range users {
			if id != userID {
				filteredUsers = append(filteredUsers, id)
			}
		}

		if len(filteredUsers) == 0 {
			if cancelFunc, ok := h.linksCancel[linkToUnTrack]; ok {
				cancelFunc()
				delete(h.linksCancel, linkToUnTrack)
			}
			delete(h.linksUsers, linkToUnTrack)
		} else {
			h.linksUsers[linkToUnTrack] = filteredUsers
		}
	}
	return nil
}
