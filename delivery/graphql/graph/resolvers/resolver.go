package resolvers

import (
	"mgufrone.dev/job-alerts/internal/usecases/job"
	"mgufrone.dev/job-alerts/internal/usecases/notification"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	uc *notification.UseCase
	jc *job.UseCase
}

func NewResolver(uc *notification.UseCase, jc *job.UseCase) *Resolver {
	return &Resolver{uc: uc, jc: jc}
}
