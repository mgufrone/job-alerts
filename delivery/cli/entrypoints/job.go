package entrypoints

import (
	"context"
	"github.com/urfave/cli"
	"mgufrone.dev/job-alerts/delivery/cli/handlers"
	"mgufrone.dev/job-alerts/internal/services/upwork"
	"mgufrone.dev/job-alerts/internal/services/weworkremotely"
	"mgufrone.dev/job-alerts/pkg/worker"
)

func NewJobEntrypoint(hndl *handlers.Job) cli.Command {
	return cli.Command{
		Name:      "job",
		Aliases:   []string{"j"},
		Usage:     "fetch jobs from available sources",
		ArgsUsage: "[upwork, weworkremotely]",
		Action: func(cCtx *cli.Context) error {
			workers := map[string]worker.IWorker{
				"upwork":         upwork.Default(),
				"weworkremotely": weworkremotely.Default(),
			}
			source := cCtx.Args().First()
			debug := cCtx.Bool("debug")
			dryRun := cCtx.Bool("dry-run")
			return hndl.Pull(context.TODO(), handlers.PullParams{
				Worker: workers[source],
				Debug:  debug,
				DryRun: dryRun,
			})
		},
	}
}
