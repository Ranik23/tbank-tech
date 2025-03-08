package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	git "tbank/scrapper/pkg/github"
	githubUtils "tbank/scrapper/pkg/github/utils"

	"github.com/IBM/sarama"
	"github.com/google/go-github/v69/github"
)

type Client struct {
	linkToTrack		string
	kafkaProducer   sarama.AsyncProducer
	logger          *slog.Logger
	gitClient       git.GitHubClient
	latestCommitSHA string
}

func NewClient(producer sarama.AsyncProducer, logger *slog.Logger, linkToTrack string, gitClient git.GitHubClient) *Client {
	return &Client{
		kafkaProducer:   	producer,
		logger:          	logger,
		gitClient:       	gitClient,
		linkToTrack: 		linkToTrack,
		latestCommitSHA: 	"",
	}
}

func (c *Client) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Stopping tracking", slog.String("link", c.linkToTrack))
			return nil
		case <-ticker.C:
			if err := c.checkForNewCommits(); err != nil {
				return err
			}
		}
	}
}

func (c *Client) checkForNewCommits() error {
	commit, err := c.getLatestCommit()
	if err != nil {
		c.logger.Error("Failed to get latest commit", slog.String("error", err.Error()))
		return err
	}

	if commit == nil || commit.GetSHA() == c.latestCommitSHA {
		return err
	}

	c.latestCommitSHA = commit.GetSHA()

	if err := c.publishCommit(commit); err != nil {
		return err
	}

	return nil
}

func (c *Client) getLatestCommit() (*github.RepositoryCommit, error) {
	owner, repo, err := githubUtils.GetLinkParams(c.linkToTrack)
	if err != nil {
		c.logger.Error("Failed to parse link", slog.String("error", err.Error()))
		return nil, err
	}
	commit, _, err := c.gitClient.LatestCommit(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}
	return commit, nil
}

func (c *Client) publishCommit(commit *github.RepositoryCommit) error {
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
		return err
	}

	owner, repo, err := githubUtils.GetLinkParams(c.linkToTrack)
	if err != nil {
		return err
	}

	topicName := fmt.Sprintf("%s_%s", owner, repo)

	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.ByteEncoder(messageValue),
	}

	c.kafkaProducer.Input() <- msg

	return nil
}