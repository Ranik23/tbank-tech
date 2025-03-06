package hub

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	git "tbank/scrapper/pkg/github"
	githubUtils"tbank/scrapper/pkg/github/utils"

	"github.com/IBM/sarama"
	"github.com/google/go-github/v69/github"
)
type Client struct {
	kafkaProducer   sarama.AsyncProducer
	logger          *slog.Logger
	gitClient       git.GitHubClient
	topic           string
	latestCommitSHA string
}

func NewClient(producer sarama.AsyncProducer, logger *slog.Logger, topic string, gitClient git.GitHubClient) *Client {
	return &Client{
		kafkaProducer:   producer,
		logger:          logger,
		gitClient:       gitClient,
		topic:           topic,
		latestCommitSHA: "",
	}
}

func (c *Client) Run(ctx context.Context, link string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Stopping tracking", slog.String("link", link))
			return
		case <-ticker.C:
			c.checkForNewCommits(link)
		}
	}
}

func (c *Client) checkForNewCommits(link string) {
	owner, repo, err := githubUtils.GetLinkParams(link)
	if err != nil {
		c.logger.Error("Failed to parse link", slog.String("error", err.Error()))
		return
	}

	commit, err := c.getLatestCommit(owner, repo)
	if err != nil {
		c.logger.Error("Failed to get latest commit", slog.String("error", err.Error()))
		return
	}

	if commit == nil || commit.GetSHA() == c.latestCommitSHA {
		return
	}

	c.latestCommitSHA = commit.GetSHA()
	c.publishCommit(commit)
}

func (c *Client) getLatestCommit(owner, repo string) (*github.RepositoryCommit, error) {
	commits, _, err := c.gitClient.ListCommits(context.Background(), owner, repo, &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 1},
	})
	if err != nil || len(commits) == 0 {
		return nil, err
	}
	return commits[0], nil
}

func (c *Client) publishCommit(commit *github.RepositoryCommit) {
	data := struct {
		SHA     string `json:"sha"`
		Message string `json:"message"`
	}{
		SHA:     commit.GetSHA(),
		Message: commit.Commit.GetMessage(),
	}

	messageValue, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("Failed to marshal commit data", slog.String("error", err.Error()))
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: c.topic,
		Value: sarama.ByteEncoder(messageValue),
	}

	c.kafkaProducer.Input() <- msg
}