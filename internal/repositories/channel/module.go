package channel

import (
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/pkg/db"
)

var RepoModule = fx.Provide(
	fx.Annotate(
		func(db *db.Instance) channel.QueryResolver {
			return func() channel.IQueryRepository {
				return NewQueryRepository(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
	fx.Annotate(
		func(db *db.Instance) channel.CommandResolver {
			return func() channel.ICommandRepository {
				return NewCommand(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
)
