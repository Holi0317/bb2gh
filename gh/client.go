package gh

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

type Client struct {
	token string
	owner string
	repo  string
}

func New(token, repository string) (*Client, error) {
	s := strings.Split(repository, "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("Malformed repository name. %s", repository)
	}

	client := &Client{
		token: token,
		owner: s[0],
		repo:  s[1],
	}

	return client, nil
}

func (c *Client) get(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
