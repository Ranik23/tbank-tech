package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)


type Client struct {
	Token 	string
	client 	*http.Client
	latestCommits map[string]string
	mut sync.Mutex
}


func NewClient(token string) *Client {
	return &Client{
		Token: token,
		client: &http.Client{},
		latestCommits: make(map[string]string),
	}
}

type Commit struct {
	SHA string `json:"sha"`
}

func (c *Client) GetTheLatestCommit(ctx context.Context, url string) (*Commit, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer " + c.Token)

	req.Header.Set("Accept", "application/json")


	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	c.mut.Lock()
	defer c.mut.Unlock()

	if c.latestCommits[url] != commits[0].SHA {
		c.latestCommits[url] = commits[0].SHA
		return &commits[0], nil
	}

	return nil, fmt.Errorf("no new commits")
}





