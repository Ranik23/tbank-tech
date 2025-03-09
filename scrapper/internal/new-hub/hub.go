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

type Pair [2]string

type Hub struct {
	gitClient       git.GitHubClient

	pairCancelFunc  *syncmap.SyncMap[Pair, context.CancelFunc] 
	latestCommitSHA *syncmap.SyncMap[string, string]


	commitChan      chan *github.RepositoryCommit
	stopCh          chan struct{}

	
	logger          *slog.Logger
}

func NewHub(gitClient git.GitHubClient, commitChan chan *github.RepositoryCommit, logger *slog.Logger) *Hub {
	return &Hub{
		gitClient:       gitClient,
		commitChan:      commitChan,
		logger:          logger,
		latestCommitSHA: syncmap.NewSyncMap[string, string](),           
		pairCancelFunc:  syncmap.NewSyncMap[Pair, context.CancelFunc](), 
		stopCh:          make(chan struct{}),
	}
}

func (h *Hub) Stop() {
	h.logger.Info("Hub is stopped")
	defer close(h.commitChan)
	h.stopCh <- struct{}{}
}

func (h *Hub) AddLink(link string, userID uint) {
	pair := Pair{link, strconv.Itoa(int(userID))}

	owner, repo, err := utils.GetLinkParams(link)
	if err != nil {
		h.logger.Error("Wrong URL scheme", slog.String("err", err.Error()))
		return
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
	
		h.pairCancelFunc.Store(pair, cancel)
	
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				h.logger.Info("Context cancelled", slog.String("link", link), slog.Int("userID", int(userID)))
				return
			case <-ticker.C:

				timeoutCtx, cancelTimeOut := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancelTimeOut()

				commit, _, err := h.gitClient.LatestCommit(timeoutCtx, owner, repo, nil)
				if err != nil {
					h.logger.Error("Failed to fetch the latest commit", slog.String("err", err.Error()))
				}

				val, ok := h.latestCommitSHA.Load(link)

				if !ok { 
					h.latestCommitSHA.Store(link, commit.GetSHA())

					h.commitChan <- commit
				} else {
					if commit.GetSHA() != val {
						h.latestCommitSHA.Store(link, commit.GetSHA())

						h.commitChan <- commit
					}
				}
			case <-h.stopCh:
				cancel()
				h.logger.Info("Goroutine %s %d exited", link, int(userID))
				return
			}
		}
	}()
}

func (h *Hub) RemoveLink(link string, userID uint) {

	pair := Pair{link, strconv.Itoa(int(userID))}
	cancelFuncForPair, ok := h.pairCancelFunc.Load(pair)

	if !ok {
		return
	} else {
		cancelFuncForPair()
		h.pairCancelFunc.Delete(pair)
	}

}
