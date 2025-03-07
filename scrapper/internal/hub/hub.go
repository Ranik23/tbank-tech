package hub

import (
	"context"
	"log/slog"
	"sync"
	"time"
	"github.com/IBM/sarama"
	git "tbank/scrapper/pkg/github"
)

type Hub struct {
	linksCancel   		map[string]context.CancelFunc
	linksUsers    		map[string][]uint
	kafkaProducer 		sarama.AsyncProducer
	mut           		sync.Mutex
	logger        		*slog.Logger
	gitClient     		git.GitHubClient
	topicToProduceIn 	string
}

func NewHub(producer sarama.AsyncProducer, logger *slog.Logger, gitClient git.GitHubClient, topic string) *Hub {
	return &Hub{
		linksCancel:   make(map[string]context.CancelFunc),
		linksUsers:    make(map[string][]uint),
		kafkaProducer: producer,
		logger:        logger,
		gitClient:     gitClient,
		topicToProduceIn: topic,
	}
}

func (h *Hub) AddTrack(link string, userID uint) {
	h.mut.Lock()
	defer h.mut.Unlock()

	if link == "" {
		return
	}

	if _, exists := h.linksUsers[link]; !exists {

		ctx, cancel := context.WithCancel(context.Background())
		client := NewClient(h.kafkaProducer, h.logger, h.topicToProduceIn, h.gitClient)

		go client.Run(ctx, link, 10*time.Second)

		h.linksCancel[link] = cancel
	}

	h.linksUsers[link] = append(h.linksUsers[link], userID)
}

func (h *Hub) RemoveTrack(link string, userID uint) {
	h.mut.Lock()
	defer h.mut.Unlock()

	if users, exists := h.linksUsers[link]; exists {
		// Удаляем userID из списка
		filteredUsers := []uint{}
		for _, id := range users {
			if id != userID {
				filteredUsers = append(filteredUsers, id)
			}
		}

		if len(filteredUsers) == 0 {
			// Если больше нет пользователей, отменяем контекст и удаляем запись
			if cancelFunc, ok := h.linksCancel[link]; ok {
				cancelFunc()
				delete(h.linksCancel, link)
			}
			delete(h.linksUsers, link)
		} else {
			h.linksUsers[link] = filteredUsers
		}
	}
}
