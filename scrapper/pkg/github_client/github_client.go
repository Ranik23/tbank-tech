package githubclient

import (
	"context"

	"github.com/google/go-github/v69/github"
)


type GitHubClient interface {
	ListCommits(ctx context.Context, owner, repo string, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
	LatestCommit(ctx context.Context, owner string, repo string, opts *github.CommitsListOptions) (*github.RepositoryCommit, *github.Response, error) 
}

type GitHubClientImpl struct {
	client github.Client
}

func NewRealGitHubClient() GitHubClient {
	return &GitHubClientImpl{}
}

func (g *GitHubClientImpl) ListCommits(ctx context.Context, owner, repo string,
							opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	return g.client.Repositories.ListCommits(ctx, owner, repo, opts)
}

func (g *GitHubClientImpl) LatestCommit(ctx context.Context, owner string, repo string,
							opts *github.CommitsListOptions) (*github.RepositoryCommit, *github.Response, error) {
	commits, response, err := g.client.Repositories.ListCommits(ctx, owner, repo, opts)
	if err != nil {
		return nil, nil, err
	}
	return commits[0], response, nil
}