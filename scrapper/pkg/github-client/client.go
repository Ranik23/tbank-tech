package githubclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const githubAPI = "https://api.github.com/repos"

type Client struct {
	HTTP *http.Client
}

type Commit struct {
	SHA    string `json:"sha"`
	Author struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Date  string `json:"date"`
	} `json:"commit.author"`
	Message string `json:"commit.message"`
}

type Version struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Created string `json:"created_at"`
}

func NewClient() *Client {
	return &Client{
		HTTP: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetTheLatestCommit(owner, repo string) (*Commit, error) {
	url := fmt.Sprintf("%s/%s/%s/commits/main", githubAPI, owner, repo)
	resp, err := c.HTTP.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch latest commit")
	}

	var commit Commit
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

func (c *Client) GetTheLatestVersion(owner, repo string) (*Version, error) {
	url := fmt.Sprintf("%s/%s/%s/releases/latest", githubAPI, owner, repo)
	resp, err := c.HTTP.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch latest version")
	}

	var version Version
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return nil, err
	}

	return &version, nil
}
