package job

import (
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/db"
)

var RepoModule = fx.Provide(
	fx.Annotate(
		func(db *db.Instance) (job.QueryResolver, job.QueryResolver) {
			resolver := func() job.IQueryRepository {
				return NewQueryRepository(db)
			}
			return resolver, resolver
		},
		fx.ResultTags(`name:"source"`, ""),
	),
	fx.Annotate(
		func(db *db.Instance) (job.CommandResolver, job.CommandResolver) {
			resolver := func() job.ICommandRepository {
				return NewCommand(db)
			}
			return resolver, resolver
		},
		fx.ResultTags(`name:"source"`, ""),
	),
)
