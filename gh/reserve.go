package gh

import (
	"errors"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/sirupsen/logrus"
)

func (c *Client) GetMaxIssueNumber() (int, error) {
	issueOpts := &github.IssueListByRepoOptions{
		State:     "all",
		Sort:      "created",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	}
	issues, _, err := c.ghc.Issues.ListByRepo(c.ctx, c.owner, c.repo, issueOpts)
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
	prs, _, err := c.ghc.PullRequests.List(c.ctx, c.owner, c.repo, prOpts)
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

func (c *Client) CreateIssue() (*github.Issue, error) {
	issueReq := &github.IssueRequest{
		Title:  github.String("Reserved by bb2gh"),
		Labels: &[]string{"bb2gh"},
	}

	for i := 0; i < retry; i++ {
		log := logrus.WithField("retry", i)

		issue, _, err := c.ghc.Issues.Create(c.ctx, c.owner, c.repo, issueReq)
		if sleep, ok := isRateLimit(err); ok {
			log.WithField("sleep", sleep).Debug("Hit rate limit. Sleeping before retry")
			time.Sleep(sleep)
			continue
		}

		if err != nil {
			return nil, err
		}

		return issue, nil
	}

	return nil, errors.New("Rate limit, (CreateIssue)")
}

func (c *Client) CloseIssue(issue *github.Issue) error {
	log := logrus.WithField("number", issue.GetNumber())
	log.Debug("Created issue. Now marking it as closed")

	updateReq := &github.IssueRequest{
		State:       github.String("closed"),
		StateReason: github.String("completed"),
	}

	for i := 0; i < retry; i++ {
		_, _, err := c.ghc.Issues.Edit(c.ctx, c.owner, c.repo, issue.GetNumber(), updateReq)
		if sleep, ok := isRateLimit(err); ok {
			log.WithField("sleep", sleep).Debug("Hit rate limit. Sleeping before retry")
			time.Sleep(sleep)
			continue
		}

		if err != nil {
			return err
		}

		log.Info("Created issue placeholder")
		return nil
	}

	return errors.New("Rate limit (Create Issue)")
}
