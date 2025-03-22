package hub

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	git "github.com/Ranik23/tbank-tech/scrapper/pkg/github_client"
	"github.com/Ranik23/tbank-tech/scrapper/pkg/github_client/utils"
	"github.com/Ranik23/tbank-tech/scrapper/pkg/syncmap"

	"github.com/google/go-github/v69/github"
)

type Pair [2]string

type CustomCommit struct {
	Commit *github.RepositoryCommit `json:"commit"`
	UserID uint                     `json:"user_id"`
}

type Hub interface {
	Run()
	Stop()
	AddLink(link string, userID uint, token string, interval time.Duration) error
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

func NewHub(gitClient git.GitHubClient, commitChan chan CustomCommit,
	logger *slog.Logger) Hub {
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

func (h *hub) AddLink(link string, userID uint, token string, interval time.Duration) error {
	const op = "Hub.AddLink"

	pair := Pair{link, strconv.Itoa(int(userID))}
	owner, repo, err := utils.GetLinkParams(link)
	if err != nil {
		h.logError(op, "Wrong URL scheme", err)
		return fmt.Errorf("wrong url scheme")
	}

	theLatestCommit, err := h.fetchLatestCommit(owner, repo, token)
	if err != nil {
		h.logError(op, "Failed to get the latest commit", err)
		return err
	}
	h.latestCommitSHA.Store(link, theLatestCommit.GetSHA())

	ctx, cancel := context.WithCancel(context.Background())
	h.pairCancelFunc.Store(pair, cancel)

	go h.trackCommits(ctx, link, owner, repo, token, userID, interval, op)

	return nil
}

func (h *hub) fetchLatestCommit(owner, repo, token string) (*github.RepositoryCommit, error) {
	commit, _, err := h.gitClient.LatestCommit(context.Background(), owner, repo, token, nil)
	if err != nil {
		return nil, err
	}
	return commit, nil
}

func (h *hub) trackCommits(ctx context.Context, link, owner, repo, token string, userID uint, interval time.Duration, op string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			h.logInfo(op, "Context cancelled", link, userID)
			return

		case <-ticker.C:
			h.checkForNewCommit(link, owner, repo, token, userID, op)

		case <-h.stopCh:
			h.logInfo(op, "Goroutine exited", link, userID)
			return
		}
	}
}

func (h *hub) checkForNewCommit(link, owner, repo, token string, userID uint, op string) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	commit, _, err := h.gitClient.LatestCommit(timeoutCtx, owner, repo, token, nil)
	if err != nil {
		h.logError(op, "Failed to fetch the latest commit", err)
		return
	}

	val, ok := h.latestCommitSHA.Load(link)
	if !ok || commit.GetSHA() != val {
		h.latestCommitSHA.Store(link, commit.GetSHA())
		h.logger.Info(op, "New Commit!", slog.String("owner", owner), slog.String("repo", repo))
		h.commitChan <- CustomCommit{Commit: commit, UserID: userID}
	}
}


func (h *hub) logError(op, msg string, err error) {
	h.logger.Error(op, slog.String("message", msg), slog.String("err", err.Error()))
}

func (h *hub) logInfo(op, msg, link string, userID uint) {
	h.logger.Info(op, slog.String("message", msg), slog.String("link", link), slog.Int("userID", int(userID)))
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
