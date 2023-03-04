package user

import (
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/pkg/db"
)

var RepoModule = fx.Provide(
	fx.Annotate(
		func(db *db.Instance) user.QueryResolver {
			return func() user.IQueryRepository {
				return NewQueryRepository(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
	fx.Annotate(
		func(db *db.Instance) user.CommandResolver {
			return func() user.ICommandRepository {
				return NewCommand(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
)
