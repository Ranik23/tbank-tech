package github

import (
	"context"

	"github.com/google/go-github/v69/github"
)


type GitHubClient interface {
	ListCommits(ctx context.Context, owner, repo string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
}

type GitHubClientImpl struct {
	client *github.Client
}

func NewRealGitHubClient() *GitHubClientImpl {
	return &GitHubClientImpl{}
}

func (g *GitHubClientImpl) ListCommits(ctx context.Context, owner, repo string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	return g.client.Repositories.ListCommits(ctx, owner, repo, opts)
}