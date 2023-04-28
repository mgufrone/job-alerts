package helper

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
)

func GetPreloads(ctx context.Context) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string) (preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(preloads, GetNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...)
	}
	return
}

func GetEdgesPreload(ctx context.Context, fieldName string) (preloads []string) {
	ctx2 := graphql.GetOperationContext(ctx)
	fields := graphql.CollectFieldsCtx(ctx, nil)
	for _, column := range fields {
		if column.Name == fieldName {
			preloads = GetNestedPreloads(ctx2, graphql.CollectFields(ctx2, column.Selections, nil), "")
			return
		}
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
