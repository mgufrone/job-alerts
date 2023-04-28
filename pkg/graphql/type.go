package graphql

import (
	"github.com/graphql-go/graphql"
	"go.uber.org/fx"
)

type IHandler interface {
	Operation() string
	Field() *graphql.Field
}

func Handler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(IHandler)),
		fx.ResultTags(`group:"graph_handlers"`),
	)
}
