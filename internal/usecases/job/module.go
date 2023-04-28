package job

import "go.uber.org/fx"

var Module = fx.Provide(
	fx.Annotate(
		NewUseCase,
		fx.ParamTags(
			`name:"source"`,
			`name:"source"`,
			`name:"source"`,
		),
	),
)
