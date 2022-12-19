package gh

import (
	"context"

	"github.com/google/go-github/v48/github"
	"github.com/sirupsen/logrus"
)

func (c *Client) PrepareLabel(ctx context.Context) error {
	client := c.get(ctx)
	logrus.Info("Getting label on github")

	label, resp, err := client.Issues.GetLabel(ctx, c.owner, c.repo, "bb2gh")
	if resp != nil && resp.StatusCode == 404 {
		return c.createLabel(ctx)
	} else if err != nil {
		return err
	}

	logrus.WithField("id", label.GetID()).Debug("Got label from github")

	return nil
}

func (c *Client) createLabel(ctx context.Context) error {
	client := c.get(ctx)

	label := &github.Label{
		Name:        github.String("bb2gh"),
		Description: github.String("Issues crated by bb2gh migration tool"),
	}

	logrus.Info("Creating label bb2gh on github")

	_, _, err := client.Issues.CreateLabel(ctx, c.owner, c.repo, label)
	if err != nil {
		return err
	}

	logrus.WithField("id", label.GetID()).Debug("Created label bb2gh")

	return nil
}
