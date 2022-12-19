package cmd

import (
	"context"

	"github.com/holi0317/bb2gh/gh"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func reserve(ctx context.Context, g *gh.Client, to int) error {
	err := g.PrepareLabel(ctx)
	if err != nil {
		return err
	}

	issueNum, err := g.GetMaxIssueNumber(ctx)
	if err != nil {
		return err
	}

	if issueNum >= to {
		logrus.WithFields(logrus.Fields{
			"to":       to,
			"issueNum": issueNum,
		}).Info("Issue number has reached the required number. Not opening new issues.")
		return nil
	}

	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(10)

	for i := issueNum + 1; i <= to; i++ {
		group.Go(func() error {
			issue, err := g.CreateIssue(ctx)
			if err != nil {
				return err
			}

			err = g.CloseIssue(ctx, issue)
			if err != nil {
				return err
			}

			return nil
		})
	}

	err = group.Wait()
	if err != nil {
		return err
	}

	return nil
}

func reserveCmd() *cli.Command {
	var flagTo int
	var flagGithubToken string
	var flagGithubRepo string

	return &cli.Command{
		Name:  "reserve",
		Usage: "Open empty Github issues to bump up the issue number",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "to",
				Usage:       "The issue number of github issue to reserve up to (inclusive).",
				Required:    true,
				Destination: &flagTo,
			},
			&cli.StringFlag{
				Name:        "github-token",
				Required:    true,
				Usage:       "Github token for accessing the repository",
				Destination: &flagGithubToken,
				EnvVars:     []string{"GITHUB_TOKEN"},
			},
			&cli.StringFlag{
				Name:        "github-repository",
				Required:    true,
				Usage:       "Github repository, in format of owner/repo",
				Destination: &flagGithubRepo,
			},
		},
		Action: func(cCtx *cli.Context) error {
			g, err := gh.New(flagGithubToken, flagGithubRepo)
			if err != nil {
				return err
			}

			err = reserve(cCtx.Context, g, flagTo)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
