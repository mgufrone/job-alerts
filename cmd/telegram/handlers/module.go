package handlers

import (
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/pkg/event"
)

var Module = fx.Options(
	fx.Provide(
		event.AsHandler(NewHelp),
		event.AsHandler(newJobHandler),
	),
)
