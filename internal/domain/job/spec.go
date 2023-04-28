package job

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"time"
)

func WhereSource(source string) criteria.ICondition {
	if source == "" {
		return nil
	}
	return criteria.NewCondition("source", criteria.Eq, source)
}
func WhereURL(urlString string) criteria.ICondition {
	if urlString == "" {
		return nil
	}
	return criteria.NewCondition("job_url", criteria.Eq, urlString)
}
func WhereID(id uuid.UUID) criteria.ICondition {
	if id == uuid.Nil {
		return nil
	}
	return criteria.NewCondition("id", criteria.Eq, id.String())
}
func WhereInID(id []uuid.UUID) criteria.ICondition {
	if len(id) == 0 {
		return nil
	}
	return criteria.NewCondition("id", criteria.In, id)
}
func WhereCompany(id string) criteria.ICondition {
	if id == "" {
		return nil
	}
	return criteria.NewCondition("company", criteria.Eq, id)
}
func WhereRole(id string) criteria.ICondition {
	if id == "" {
		return nil
	}
	return criteria.NewCondition("role", criteria.Eq, id)
}
func WhereJobType(id string) criteria.ICondition {
	if id == "" {
		return nil
	}
	return criteria.NewCondition("job_type", criteria.Eq, id)
}
func WhereTags(vals []string) criteria.ICondition {
	if len(vals) == 0 {
		return nil
	}
	return criteria.NewCondition("name", criteria.In, vals)
}
func WhereRoleContains(val string) criteria.ICondition {
	if val == "" {
		return nil
	}
	return criteria.NewCondition("role", criteria.Like, val)
}
func WhereDescriptionContains(val string) criteria.ICondition {
	if val == "" {
		return nil
	}
	return criteria.NewCondition("description", criteria.Like, val)
}
func WhereMinSalary(val float64) criteria.ICondition {
	if val == 0 {
		return nil
	}
	return criteria.NewCondition("salary", criteria.Like, val)
}
func WhereBefore(duration time.Duration) criteria.ICondition {
	now := time.Now()
	prev := now.Add(-duration)
	return criteria.NewCondition("created_at", criteria.Between, []time.Time{prev, now})
}
