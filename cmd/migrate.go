package cmd

import "github.com/urfave/cli/v2"

func migrateCmd() *cli.Command {
	var flagFrom int
	var flagTo int

	return &cli.Command{
		Name:  "migrate",
		Usage: "Migrate given range of bitbucket PR to Github issue. The github issues must been reserved beforehand.",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "from",
				Usage:       "Starting PR number",
				Required:    true,
				Destination: &flagFrom,
			},
			&cli.IntFlag{
				Name:        "to",
				Usage:       "Ending PR number",
				Required:    true,
				Destination: &flagTo,
			},
		},
		Action: func(cCtx *cli.Context) error {
			return nil
		},
	}
}
