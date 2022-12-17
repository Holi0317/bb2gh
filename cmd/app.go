package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	return &cli.App{
		Name:  "bb2gh",
		Usage: "Migrate PR from Bitbucket to Github as issue",
		Commands: []*cli.Command{
			migrateCmd(),
			reserveCmd(),
		},
		Before: func(ctx *cli.Context) error {
			logrus.SetLevel(logrus.DebugLevel)
			return nil
		},
	}
}
