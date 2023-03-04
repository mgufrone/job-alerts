package job

import (
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/db"
)

var RepoModule = fx.Provide(
	fx.Annotate(
		func(db *db.Instance) job.QueryResolver {
			return func() job.IQueryRepository {
				return NewQueryRepository(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
	fx.Annotate(
		func(db *db.Instance) job.CommandResolver {
			return func() job.ICommandRepository {
				return NewCommand(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
)
