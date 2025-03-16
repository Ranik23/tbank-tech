package hub

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	git "tbank/scrapper/pkg/github_client"
	"tbank/scrapper/pkg/github_client/utils"
	"tbank/scrapper/pkg/syncmap"
	"time"

	"github.com/google/go-github/v69/github"
)

type Pair [2]string

type CustomCommit struct {
	Commit *github.RepositoryCommit	`json:"commit"`
	UserID uint						`json:"user_id"`
}

type Hub interface {
	Run()
	Stop()
	AddLink(link string, userID uint) 	 error
	RemoveLink(link string, userID uint) error
}

type hub struct {
	gitClient       git.GitHubClient
	pairCancelFunc  *syncmap.SyncMap[Pair, context.CancelFunc]
	latestCommitSHA *syncmap.SyncMap[string, string]
	commitChan      chan CustomCommit
	stopCh          chan struct{}
	logger          *slog.Logger
}

func NewHub(gitClient git.GitHubClient, commitChan chan CustomCommit, logger *slog.Logger) Hub{
	const op = "Hub.NewHub"
	logger.Info(op, slog.String("message", "Creating new Hub"))

	return &hub{
		gitClient:       gitClient,
		commitChan:      commitChan,
		logger:          logger,
		latestCommitSHA: syncmap.NewSyncMap[string, string](),
		pairCancelFunc:  syncmap.NewSyncMap[Pair, context.CancelFunc](),
		stopCh:          make(chan struct{}),
	}
}

func (h *hub) Run() {
	const op = "Hub.Run"
	h.logger.Info(op, slog.String("message", "Hub is running..."))

	go func() { 
		<-h.stopCh
		h.logger.Info(op, slog.String("message", "Hub is stopping..."))

		h.pairCancelFunc.Range(func(pair Pair, cancel context.CancelFunc) bool {
			cancel()
			return true
		})
		close(h.commitChan)
	}()
}

func (h *hub) AddLink(link string, userID uint) error {
	const op = "Hub.AddLink"
	pair := Pair{link, strconv.Itoa(int(userID))}
	owner, repo, err := utils.GetLinkParams(link)
	if err != nil {
		h.logger.Error(op, slog.String("message", "Wrong URL scheme"), slog.String("err", err.Error()))
		return fmt.Errorf("wrong url scheme")
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.pairCancelFunc.Store(pair, cancel)

	go func() {
		ticker := time.NewTicker(4 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				h.logger.Info(op, slog.String("message", "Context cancelled"), slog.String("link", link), slog.Int("userID", int(userID)))
				return 
			case <-ticker.C:

				timeoutCtx, cancelTimeOut := context.WithTimeout(context.Background(), 10*time.Second)

				commit, _, err := h.gitClient.LatestCommit(timeoutCtx, owner, repo, nil)
				if err != nil {
					h.logger.Error(op, slog.String("message", "Failed to fetch the latest commit"), slog.String("err", err.Error()))
					cancelTimeOut()
					continue
				}

				cancelTimeOut()

				val, ok := h.latestCommitSHA.Load(link)

				if !ok || commit.GetSHA() != val {
					h.latestCommitSHA.Store(link, commit.GetSHA())
					h.commitChan <- CustomCommit{Commit: commit, UserID: userID}
				}
			case <-h.stopCh:
				h.logger.Info(op, slog.String("message", "Goroutine exited"), slog.String("link", link), slog.Int("userID", int(userID)))
				return 
			}
		}
	}()

	return nil
}

func (h *hub) Stop() {
	const op = "Hub.Stop"
	h.logger.Info(op, slog.String("message", "Stopping Hub..."))
	close(h.stopCh)
}

func (h *hub) RemoveLink(link string, userID uint) error {
	const op = "Hub.RemoveLink"
	pair := Pair{link, strconv.Itoa(int(userID))}
	cancelFuncForPair, ok := h.pairCancelFunc.Load(pair)

	if !ok {
		return fmt.Errorf("pair doesn't exist")
	} else {
		cancelFuncForPair()
		h.pairCancelFunc.Delete(pair)
		h.logger.Info(op, slog.String("message", "Removed link"), slog.String("link", link), slog.Int("userID", int(userID)))
		return nil
	}
}
