package hub

import (
	"context"
	"net/http"
	"time"
)



type Client struct {
	client *http.Client
}


func NewClient() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}

func (c *Client) Run(ctx context.Context, link string, timer time.Duration) {


	ticker := time.NewTicker(timer)
	defer ticker.Stop()
	

	for {
		select {
		case <- ctx.Done():


	
		case <- ticker.C:



		}
	}

}