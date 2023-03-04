package notification

import (
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/pkg/db"
)

var RepoModule = fx.Provide(
	fx.Annotate(
		func(db *db.Instance) notification.QueryResolver {
			return func() notification.IQueryRepository {
				return NewQueryRepository(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
	fx.Annotate(
		func(db *db.Instance) notification.CommandResolver {
			return func() notification.ICommandRepository {
				return NewCommand(db)
			}
		},
		fx.ResultTags(`name:"source"`),
	),
)
