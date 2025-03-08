package newhub

import (
	"context"
	"log/slog"
	"strconv"
	git "tbank/scrapper/pkg/github"
	"tbank/scrapper/pkg/github/utils"
	"tbank/scrapper/pkg/syncmap"
	"time"

	"github.com/google/go-github/v69/github"
)

type Pair [2]string // {link, userID}


type Hub struct {
	gitClient 		git.GitHubClient
	pairCancelFunc 	*syncmap.SyncMap[Pair, context.CancelFunc] //map[Pair]context.CancelFunc
	logger			*slog.Logger
	latestCommitSHA *syncmap.SyncMap[string, string]           //map[string]string // default: ""
	commitChan		chan *github.RepositoryCommit
	stopCh			chan struct{}
}


func NewHub(gitClient git.GitHubClient, commitChan chan *github.RepositoryCommit, logger *slog.Logger) *Hub {
	return &Hub{
		gitClient: gitClient,
		commitChan: commitChan,
		logger: logger,
		latestCommitSHA: syncmap.NewSyncMap[string, string](), //make(map[string]string),
		pairCancelFunc: syncmap.NewSyncMap[Pair, context.CancelFunc](), //make(map[Pair]context.CancelFunc),
		stopCh: make(chan struct{}),
	}
}



func (h *Hub) Run() {
	go func() {
		for {
			select {
			case <-h.commitChan:
				
			case <-h.stopCh:
				return
			}
		}
	}()
}


func (h *Hub) Result() <- chan *github.RepositoryCommit {
	return h.commitChan
} 

func (h *Hub) Stop() {
	h.logger.Info("Hub is stopped")
	h.stopCh <- struct{}{}
}


func (h *Hub) AddLink(link string, userID uint) {
	pair := Pair{link, strconv.Itoa(int(userID))}

	go func() {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		h.pairCancelFunc.Store(pair, cancel)

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				h.logger.Info("Context cancelled", slog.String("link", link), slog.Int("userID", int(userID)))
				return
			case <- ticker.C:
				owner, repo, err := utils.GetLinkParams(link)
				if err != nil {
					h.logger.Error("Wrong URL scheme", slog.String("err", err.Error()))
					return
				}

				timeoutCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
				defer cancel()
				
				commit, _, err := h.gitClient.LatestCommit(timeoutCtx, owner, repo, nil)
				if err != nil {
					h.logger.Error("Failed to fetch the latest commit", slog.String("err", err.Error()))
				}

				val, ok := h.latestCommitSHA.Load(link)

				if !ok {
					h.latestCommitSHA.Store(link, commit.GetSHA())
				} else {
					if commit.GetSHA() != val {
						h.latestCommitSHA.Store(link, commit.GetSHA())

						h.commitChan <- commit
					}
				}
			}
		}
	}()

}

func (h *Hub) RemoveLink(link string, userID uint) {

	pair := Pair{link, strconv.Itoa(int(userID))}
	cancelFuncForPair, ok := h.pairCancelFunc.Load(pair)

	if !ok { // следовательно тако пары нет
		return
	} else {
		// просто отменяем ту горутину которая была запущена
		cancelFuncForPair()
	}

}