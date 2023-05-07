package helpers

import (
	"go.uber.org/fx"
)

func AsCommand(f any) any {
	return fx.Annotate(
		f,
		fx.ResultTags(`group:"commands"`),
	)
}
