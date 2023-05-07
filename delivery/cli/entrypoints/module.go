package entrypoints

import (
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/delivery/cli/helpers"
)

var Module = fx.Provide(
	helpers.AsCommand(NewJobEntrypoint),
	helpers.AsCommand(Migration),
)
