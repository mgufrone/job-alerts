package entrypoints

import (
	"context"
	"mgufrone.dev/job-alerts/delivery/cli/app"
	"mgufrone.dev/job-alerts/delivery/cli/handlers"
	"mgufrone.dev/job-alerts/delivery/cli/helpers"
	"mgufrone.dev/job-alerts/internal/services/upwork"
	"mgufrone.dev/job-alerts/internal/services/weworkremotely"
	"mgufrone.dev/job-alerts/pkg/worker"
)

type JobCmd struct {
	helpers.Context
	Source string `arg:"" name:"source" help:"pull jobs from available sources: upwork, weworkremotely"`
}

func (j *JobCmd) Run(console *app.Kernel) error {
	workers := map[string]worker.IWorker{
		"upwork":         upwork.Default(),
		"weworkremotely": weworkremotely.Default(),
	}
	return console.Job.Pull(context.TODO(), handlers.PullParams{
		Worker: workers[j.Source],
		Debug:  j.Debug,
		DryRun: j.DryRun,
	})
}
