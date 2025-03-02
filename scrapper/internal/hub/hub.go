package hub

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

type Hub struct {
	linksCancel map[string]context.CancelFunc
	linksCouns map[string]int
	producer sarama.AsyncProducer
}


func NewHub(producer sarama.AsyncProducer) *Hub {
	return &Hub{
		linksCancel: make(map[string]context.CancelFunc),
		linksCouns: make(map[string]int),
		producer : producer,
	}
}


func (h *Hub) AddTrack(link string) { 

	_, ok := h.linksCancel[link]
	if !ok {

		h.linksCouns[link] = 1

		ctx, cancel := context.WithCancel(context.Background())
		
		client := NewClient()

		go client.Run(ctx, link, 10 * time.Second)

		h.linksCancel[link] = cancel
	}else {

		h.linksCouns[link] += 1
	}
}

func (h *Hub) RemoveTrack(link string) {
	
}