package entrypoints

import (
	"github.com/urfave/cli"
	"mgufrone.dev/job-alerts/delivery/cli/handlers"
)

func Migration(mgr *handlers.Migration) cli.Command {
	return cli.Command{
		Name:  "migration",
		Usage: "Run db migration operation",
		Subcommands: []cli.Command{
			{
				Name:      "create",
				Usage:     "create new migration from predefined template",
				ArgsUsage: "<migration name>",
				Action: func(cCtx *cli.Context) error {
					return mgr.Create(cCtx.Args().First())
				},
			},
			{
				Name:  "status",
				Usage: "check migration status",
				Action: func(cCtx *cli.Context) error {
					return mgr.Status()
				},
			},
			{
				Name:  "up",
				Usage: "push & apply new migration",
				Action: func(cCtx *cli.Context) error {
					dryRun := cCtx.GlobalBool("dry-run")
					return mgr.Up(handlers.MigrationContext{DryRun: dryRun})
				},
			},
			{
				Name:  "down",
				Usage: "undo up migration",
				Action: func(cCtx *cli.Context) error {
					dryRun := cCtx.GlobalBool("dry-run")
					return mgr.Down(handlers.MigrationContext{DryRun: dryRun})
				},
			},
		},
	}
}
