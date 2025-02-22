package hub

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

type Client struct {
	client 			*github.Client
	token 			string
	latestCommit 	*github.RepositoryCommit
	updateCh 		chan *github.RepositoryCommit
	
}


func NewClient(token string, updateCh chan *github.RepositoryCommit) *Client {
	return &Client{
		client: 		github.NewClient(nil),
		token: 			token,
		updateCh: 		updateCh,
		latestCommit: 	nil,
	}
}


func (c *Client) Search(ctx context.Context, url string) error {
	owner, repo, err := c.parseGitHubURL(url)
	if err != nil {
		return fmt.Errorf("invalid url format")
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ctx.Done():
			slog.Info("context cancelled - hub closed")
			return nil
		case <- ticker.C:
			commits, _, err := c.client.Repositories.ListCommits(ctx, owner, repo, nil)
			if err != nil {
				slog.Error("failed to fetch commits", "error", err)
				continue
			}

			var newCommitsAfterTheLatest []*github.RepositoryCommit

			for i := len(commits) - 1; i >= 0; i-- {
				if c.latestCommit != nil && commits[i].GetSHA() != c.latestCommit.GetSHA() {
					newCommitsAfterTheLatest = append(newCommitsAfterTheLatest, commits[i])
				} else {
					break
				}
			}

			if len(newCommitsAfterTheLatest) > 0 {
				c.latestCommit = newCommitsAfterTheLatest[0]
			}
			
			for _, commit := range newCommitsAfterTheLatest {
				c.updateCh <- commit
			}
		}
	}
}


func (c *Client) parseGitHubURL(repoURL string) (owner, repo string, err error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("некорректный URL репозитория")
	}

	return parts[0], parts[1], nil
}