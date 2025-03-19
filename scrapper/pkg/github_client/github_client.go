package githubclient

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/go-github/v69/github"
)

type GitHubClient interface {
	ListCommits(ctx context.Context, owner, repo string, token string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
	LatestCommit(ctx context.Context, owner string, repo string, token string, opts *github.CommitsListOptions) (*github.RepositoryCommit, *github.Response, error)
}

type GitHubClientImpl struct {
	client *github.Client
	logger *slog.Logger
}

func NewRealGitHubClient(logger *slog.Logger) GitHubClient {
	return &GitHubClientImpl{
		client: github.NewClient(http.DefaultClient),
		logger: logger,
	}
}

func (g *GitHubClientImpl) ListCommits(ctx context.Context, owner, repo string, token string,
	opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {

	g.logger.Info("Fetching commits", slog.String("owner", owner), slog.String("repo", repo))

	req, err := g.client.NewRequest("GET", "repos/"+owner+"/"+repo+"/commits", opts)
	if err != nil {
		g.logger.Error("Failed to create request", slog.String("error", err.Error()))
		return nil, nil, err
	}
	req.Header.Set("Authorization", "token "+token)

	g.logger.Info("Token - ", slog.String("token", token))

	var commits []*github.RepositoryCommit
	resp, err := g.client.Do(ctx, req, &commits)
	if err != nil {
		g.logger.Error("Failed to fetch commits", slog.String("error", err.Error()))
		return nil, resp, err
	}

	g.logger.Info("Fetched commits", slog.Int("count", len(commits)))
	return commits, resp, nil
}

func (g *GitHubClientImpl) LatestCommit(ctx context.Context, owner string, repo string, token string,
	opts *github.CommitsListOptions) (*github.RepositoryCommit, *github.Response, error) {

	g.logger.Info("Fetching latest commit", slog.String("owner", owner), slog.String("repo", repo))

	req, err := g.client.NewRequest("GET", "repos/"+owner+"/"+repo+"/commits", opts)
	if err != nil {
		g.logger.Error("Failed to create request", slog.String("error", err.Error()))
		return nil, nil, err
	}
	req.Header.Set("Authorization", "token "+token)

	g.logger.Info("Token - ", slog.String("token", token))
	
	var commits []*github.RepositoryCommit
	resp, err := g.client.Do(ctx, req, &commits)
	if err != nil {
		g.logger.Error("Failed to fetch latest commit", slog.String("error", err.Error()))
		return nil, resp, err
	}

	if len(commits) == 0 {
		g.logger.Warn("No commits found")
		return nil, resp, nil
	}

	g.logger.Info("Fetched latest commit", slog.String("sha", commits[0].GetSHA()))
	return commits[0], resp, nil
}
