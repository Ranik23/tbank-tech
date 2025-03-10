package hub

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

type CustomCommit struct {
	Commit *github.RepositoryCommit
	UserID uint
}

type Hub struct {
	gitClient       git.GitHubClient
	pairCancelFunc  *syncmap.SyncMap[Pair, context.CancelFunc]
	latestCommitSHA *syncmap.SyncMap[string, string]
	commitChan      chan CustomCommit
	stopCh          chan struct{}
	logger          *slog.Logger
}

func NewHub(gitClient git.GitHubClient, commitChan chan CustomCommit, logger *slog.Logger) *Hub {
	return &Hub{
		gitClient:       gitClient,
		commitChan:      commitChan,
		logger:          logger,
		latestCommitSHA: syncmap.NewSyncMap[string, string](),
		pairCancelFunc:  syncmap.NewSyncMap[Pair, context.CancelFunc](),
		stopCh:          make(chan struct{}),
	}
}

func (h *Hub) AddLink(link string, userID uint) {
	pair := Pair{link, strconv.Itoa(int(userID))}
	owner, repo, err := utils.GetLinkParams(link)
	if err != nil {
		h.logger.Error("Wrong URL scheme", slog.String("err", err.Error()))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.pairCancelFunc.Store(pair, cancel)

	go func() {
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
					continue
				}

				val, ok := h.latestCommitSHA.Load(link)

				if !ok || commit.GetSHA() != val {
					h.latestCommitSHA.Store(link, commit.GetSHA())
					h.commitChan <- CustomCommit{Commit: commit, UserID: userID}
				}
			case <-h.stopCh:
				cancel()
				h.logger.Info("Goroutine exited", slog.String("link", link), slog.Int("userID", int(userID)))
				return
			}
		}
	}()
}


func (h *Hub) Stop() {
	h.logger.Info("Hub is stopped")
	defer close(h.commitChan)
	h.stopCh <- struct{}{}
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
