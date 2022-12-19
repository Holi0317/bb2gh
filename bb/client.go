package bb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PullRequestNotFound struct{}

func (e *PullRequestNotFound) Error() string {
	return "the repository or pull request does not exist"
}

type Client struct {
	client http.Client

	user     string
	password string

	workspace string
	reposlug  string
}

func New(user, password, repo string) (*Client, error) {
	split := strings.Split(repo, "/")
	if len(split) != 2 {
		return nil, fmt.Errorf("Is not in correct bitbucket repo format %v", repo)
	}

	client := &Client{
		client: http.Client{},

		user:     user,
		password: password,

		workspace: split[0],
		reposlug:  split[1],
	}

	return client, nil
}

func (b *Client) GetPR(ctx context.Context, prid int) (*PullRequest, error) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests/%d", b.workspace, b.reposlug, prid)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(b.user, b.password)

	req.Header.Add("accept", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("401 Unauthorized")
	}

	if resp.StatusCode == 404 {
		return nil, &PullRequestNotFound{}
	}

	var pr PullRequest
	reader := json.NewDecoder(resp.Body)

	err = reader.Decode(&pr)
	if err != nil {
		return nil, err
	}

	pr.Summary.check(prid)
	pr.Rendered.Description.check(prid)
	pr.Rendered.Title.check(prid)
	pr.Rendered.Reason.check(prid)

	return &pr, nil
}
