package gh

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

type Client struct {
	ctx context.Context

	token string
	owner string
	repo  string
	ghc   *github.Client
}

func New(ctx context.Context, token, repository string) (*Client, error) {
	s := strings.Split(repository, "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("Malformed repository name. %s", repository)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	ghc := github.NewClient(tc)

	client := &Client{
		ctx:   ctx,
		token: token,
		owner: s[0],
		repo:  s[1],
		ghc:   ghc,
	}

	return client, nil
}
