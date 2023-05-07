package handlers

import "go.uber.org/fx"

var Module = fx.Provide(
	NewJob,
	NewMigration,
)
