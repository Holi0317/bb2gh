package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/holi0317/bb2gh/bb"
	"github.com/holi0317/bb2gh/gh"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func singleMigrate(ctx context.Context, g *gh.Client, b *bb.Client, number int) error {
	log := logrus.WithField("number", number)

	log.Info("Migrate PR")
	src, err := b.GetPR(ctx, number)
	if err != nil {
		return err
	}

	log.WithField("title", src.Title).Debug("Got PR info from bitbucket. Writing to github")
	err = g.UpdateIssue(ctx, number, src)
	if err != nil {
		return err
	}

	log.Info("Migrated PR successfully")

	return nil
}

func migrate(ctx context.Context, g *gh.Client, b *bb.Client, numbers []int) error {
	group, ctx := errgroup.WithContext(ctx)
	group.SetLimit(10)

	for _, num := range numbers {
		num := num
		group.Go(func() error {
			err := singleMigrate(ctx, g, b, num)
			if err != nil {
				logrus.WithField("num", num).WithError(err).Warn("Failed to migrate issue")
			}

			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		return err
	}

	return nil
}

func migrateCmd() *cli.Command {
	var flagGithubToken string
	var flagGithubRepo string
	var flagBBUser string
	var flagBBPassword string
	var flagBBRepo string

	return &cli.Command{
		Name:      "migrate",
		Usage:     "Migrate given range of bitbucket PR to Github issue. The github issues must been reserved beforehand.",
		ArgsUsage: "List of PR numbers for migrate. Use bash sequence expression {1..100} to provide a range",
		Flags: []cli.Flag{
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
				Usage:       "Destination Github repository, in format of owner/repo",
				Destination: &flagGithubRepo,
			},
			&cli.StringFlag{
				Name:        "bitbucket-repository",
				Required:    true,
				Usage:       "Source Bitbucket repository, in format of workspace/reposlug",
				Destination: &flagBBRepo,
			},
			&cli.StringFlag{
				Name:        "bitbucket-user",
				Required:    true,
				Usage:       "Bitbucket login username (not email!)",
				Destination: &flagBBUser,
				EnvVars:     []string{"BITBUCKET_USER"},
			},
			&cli.StringFlag{
				Name:        "bitbucket-password",
				Required:    true,
				Usage:       "Bitbucket app password",
				Destination: &flagBBPassword,
				EnvVars:     []string{"BITBUCKET_PASSWORD"},
			},
		},
		Action: func(cCtx *cli.Context) error {
			args := cCtx.Args().Slice()
			numbers := make([]int, len(args))
			for i, num := range args {
				parsed, err := strconv.Atoi(num)
				if err != nil {
					return fmt.Errorf("Given argument is not an integer %s: %v", num, err)
				}

				numbers[i] = parsed
			}

			g, err := gh.New(flagGithubToken, flagGithubRepo)
			if err != nil {
				return err
			}

			b, err := bb.New(flagBBUser, flagBBPassword, flagBBRepo)
			if err != nil {
				return err
			}

			err = migrate(cCtx.Context, g, b, numbers)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
