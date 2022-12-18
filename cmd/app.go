package cmd

import (
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
	}
}
