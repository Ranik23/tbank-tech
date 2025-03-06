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
	linksCount    		map[string]int
	kafkaProducer 		sarama.AsyncProducer
	mut           		sync.Mutex
	logger        		*slog.Logger
	gitClient     		git.GitHubClient
	topicToProduceIn 	string
}

func NewHub(producer sarama.AsyncProducer, logger *slog.Logger, gitClient git.GitHubClient, topic string) *Hub {
	return &Hub{
		linksCancel:   make(map[string]context.CancelFunc),
		linksCount:    make(map[string]int),
		kafkaProducer: producer,
		logger:        logger,
		gitClient:     gitClient,
		topicToProduceIn: topic,
	}
}

func (h *Hub) AddTrack(link string) {
	h.mut.Lock()
	defer h.mut.Unlock()

	if link == "" {
		return
	}

	if count, exists := h.linksCount[link]; !exists {
		h.linksCount[link] = 1

		ctx, cancel := context.WithCancel(context.Background())
		client := NewClient(h.kafkaProducer, h.logger, h.topicToProduceIn, h.gitClient)

		go client.Run(ctx, link, 10*time.Second)

		h.linksCancel[link] = cancel
	} else {
		h.linksCount[link] = count + 1
	}
}

func (h *Hub) RemoveTrack(link string) {
	h.mut.Lock()
	defer h.mut.Unlock()

	if count, exists := h.linksCount[link]; exists {
		if count > 1 {
			h.linksCount[link] = count - 1
		} else {
			if cancelFunc, ok := h.linksCancel[link]; ok {
				cancelFunc()
				delete(h.linksCancel, link)
			}
			delete(h.linksCount, link)
		}
	}
}