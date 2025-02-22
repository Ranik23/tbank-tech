package hub

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"tbank/bot/api/proto/gen"
	"tbank/scrapper/config"
	"tbank/scrapper/internal/db/models"
	"time"

	"github.com/google/go-github/github"
	"google.golang.org/grpc"
)

// TODO: РАЗОБРАТЬСЯ С КОНТЕКСТОМ
type Hub struct {
	linkChats 			map[models.Link][]int64
	linkCancelFunc		map[models.Link]context.CancelFunc

	grpcBotClient 		gen.BotClient
	config 				*config.Config
	updatesCh			chan *github.RepositoryCommit
	stopCh				chan struct{}
	contextWithCancel	context.Context
	cancelFunc			context.CancelFunc

	token 				string
	mu 					sync.Mutex
}


func NewHub(cfg *config.Config) (*Hub, error) {

	host := cfg.Bot.Host
	port := cfg.Bot.Port

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	grpcBotClient := gen.NewBotClient(conn)

	ctx, cancel := context.WithCancel(context.Background())

	return &Hub{
		linkChats: make(map[models.Link][]int64),
		grpcBotClient: grpcBotClient,
		config: cfg,
		updatesCh: make(chan *github.RepositoryCommit),
		stopCh: make(chan struct{}),
		contextWithCancel: ctx,
		cancelFunc: cancel,
	}, nil
}


func (s *Hub) AddTrack(link models.Link, chatID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.linkChats[link]; !exists {

		s.linkChats[link] = []int64{chatID}

		c := NewClient(s.token, s.updatesCh)

		contextWithCancel, cancel := context.WithCancel(context.Background())
		s.linkCancelFunc[link] = cancel

		go c.Search(contextWithCancel, link)

	} else {
		if !slices.Contains(s.linkChats[link], chatID) {
			s.linkChats[link] = append(s.linkChats[link], chatID)
		}
	}
}


func (s *Hub) RemoveTrack(link models.Link, chatID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ids, exists := s.linkChats[link]
	if !exists {
		return fmt.Errorf("no link found")
	}

	if len(ids) == 1 && ids[0] == chatID {
		if cancelFunc, exists := s.linkCancelFunc[link]; exists {
			cancelFunc() 
			delete(s.linkCancelFunc, link)
		}
		delete(s.linkChats, link)
		return nil
	}

	newIDS := ids[:0]
	for _, id := range ids {
		if id != chatID {
			newIDS = append(newIDS, id)
		}
	}

	if len(newIDS) == 0 {
		if cancelFunc, ok := s.linkCancelFunc[link]; ok {
			cancelFunc()
			delete(s.linkCancelFunc, link)
		}
		delete(s.linkChats, link)
	} else {
		s.linkChats[link] = newIDS
	}

	return nil
}
  


func (s *Hub) Start() error {
	go func() {
		for {
			select {
			case commit := <-s.updatesCh:
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				update := &gen.CommitUpdate{
					Url: 		commit.GetURL(),
					Sha:  		commit.GetSHA(),
					Author: 	commit.Author.GetName(),
					Message: 	commit.Commit.GetMessage(),
					Timestamp: 	commit.Commit.GetAuthor().GetDate().String(),
				}

				_, err := s.grpcBotClient.SendCommitUpdate(ctx, update)
				if err != nil {
					slog.Error("failed to send update", slog.String("error", err.Error()))
				}
			case <-s.stopCh:
				slog.Info("Stopping updates processing")
				return
			}
		}
	}()
	return nil
}


func (s *Hub) Stop() error {
	for _, cancel := range s.linkCancelFunc {
		cancel()
	}
	s.stopCh <- struct{}{}
	close(s.updatesCh)
	return nil
}

