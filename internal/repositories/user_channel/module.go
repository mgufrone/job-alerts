package user_channel

import (
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/internal/domain/user_channel"
	"mgufrone.dev/job-alerts/pkg/db"
)

var RepoModule = fx.Provide(
	fx.Annotate(
		func(db *db.Instance) user_channel.QueryResolver {
			return func() user_channel.IQueryRepository {
				return NewQueryRepository(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
	fx.Annotate(
		func(db *db.Instance) user_channel.CommandResolver {
			return func() user_channel.ICommandRepository {
				return NewCommand(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
)
