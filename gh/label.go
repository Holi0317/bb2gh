package gh

import (
	"github.com/google/go-github/v48/github"
	"github.com/sirupsen/logrus"
)

func (c *Client) PrepareLabel() error {
	logrus.Info("Getting label on github")

	label, resp, err := c.ghc.Issues.GetLabel(c.ctx, c.owner, c.repo, "bb2gh")
	if resp != nil && resp.StatusCode == 404 {
		return c.createLabel()
	} else if err != nil {
		return err
	}

	logrus.WithField("id", label.GetID()).Debug("Got label from github")

	return nil
}

func (c *Client) createLabel() error {
	label := &github.Label{
		Name:        github.String("bb2gh"),
		Description: github.String("Issues crated by bb2gh migration tool"),
	}

	logrus.Info("Creating label bb2gh on github")

	_, _, err := c.ghc.Issues.CreateLabel(c.ctx, c.owner, c.repo, label)
	if err != nil {
		return err
	}

	logrus.WithField("id", label.GetID()).Debug("Created label bb2gh")

	return nil
}
