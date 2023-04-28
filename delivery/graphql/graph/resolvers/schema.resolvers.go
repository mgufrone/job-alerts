package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.25

import (
	"context"
	"fmt"

	"mgufrone.dev/job-alerts/delivery/graphql/graph/generated"
)

// Ping is the resolver for the ping field.
func (r *mutationResolver) Ping(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: Ping - ping"))
}

// Version is the resolver for the version field.
func (r *queryResolver) Version(ctx context.Context) (string, error) {
	panic(fmt.Errorf("not implemented: Version - version"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }