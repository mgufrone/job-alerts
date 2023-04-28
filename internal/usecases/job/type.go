package job

import (
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user"
)

type UseCase struct {
	jbQuery job.QueryResolver
	jbCmd   job.CommandResolver
	ntQuery notification.QueryResolver
}

func NewUseCase(jbQuery job.QueryResolver, jbCmd job.CommandResolver, ntQuery notification.QueryResolver) *UseCase {
	return &UseCase{jbQuery: jbQuery, jbCmd: jbCmd, ntQuery: ntQuery}
}

type Pagination struct {
	PerPage int
	Page    int
}

type Query struct {
	Fields []string
}

type ListInput struct {
	Keyword string
	Skills  []string
	User    *user.Entity
	IsRead  bool
	Pagination
	Query
}
