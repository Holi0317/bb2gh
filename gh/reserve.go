package gh

import (
	"context"

	"github.com/google/go-github/v48/github"
	"github.com/sirupsen/logrus"
)

func (c *Client) GetMaxIssueNumber(ctx context.Context) (int, error) {
	client := c.get(ctx)

	issueOpts := &github.IssueListByRepoOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	}
	issues, _, err := client.Issues.ListByRepo(ctx, c.owner, c.repo, issueOpts)
	if err != nil {
		return -1, err
	}

	maxIssue := 0
	for _, issue := range issues {
		maxIssue = issue.GetNumber()
	}

	prOpts := &github.PullRequestListOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	}
	prs, _, err := client.PullRequests.List(ctx, c.owner, c.repo, prOpts)
	if err != nil {
		return -1, err
	}

	maxPr := 0
	for _, pr := range prs {
		maxPr = pr.GetNumber()
	}

	logrus.WithFields(logrus.Fields{
		"maxIssue": maxIssue,
		"maxPr":    maxPr,
	}).Debug("Got issue numbers")

	if maxIssue > maxPr {
		return maxIssue, nil
	}

	return maxPr, nil
}

func (c *Client) CreateIssue(ctx context.Context) (*github.Issue, error) {
	client := c.get(ctx)

	issueReq := &github.IssueRequest{
		Title:  github.String("Reserved by bb2gh"),
		Labels: &[]string{"bb2gh"},
	}

	issue, _, err := client.Issues.Create(ctx, c.owner, c.repo, issueReq)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (c *Client) CloseIssue(ctx context.Context, issue *github.Issue) error {
	client := c.get(ctx)

	log := logrus.WithField("number", issue.GetNumber())
	log.Debug("Created issue. Now marking it as closed")

	updateReq := &github.IssueRequest{
		State:       github.String("closed"),
		StateReason: github.String("completed"),
	}

	_, _, err := client.Issues.Edit(ctx, c.owner, c.repo, issue.GetNumber(), updateReq)

	if err != nil {
		return err
	}

	return nil
}
