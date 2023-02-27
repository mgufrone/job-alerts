package job

import (
	"mgufrone.dev/job-alerts/pkg/db"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

func TagCriteria() criteria.ICriteriaBuilder {
	return db.NewCriteriaBuilder("Tags")
}
