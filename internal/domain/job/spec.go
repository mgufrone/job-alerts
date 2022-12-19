package job

import (
	criteria2 "mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

func WhereSource(source string) criteria2.ICondition {
	if source == "" {
		return nil
	}
	return criteria2.NewCondition("source", criteria2.Eq, source)
}
func WhereURL(urlString string) criteria2.ICondition {
	if urlString == "" {
		return nil
	}
	return criteria2.NewCondition("job_url", criteria2.Eq, urlString)
}
func WhereID(id string) criteria2.ICondition {
	if id == "" {
		return nil
	}
	return criteria2.NewCondition("id", criteria2.Eq, id)
}
func WhereInID(id []string) criteria2.ICondition {
	if len(id) == 0 {
		return nil
	}
	return criteria2.NewCondition("id", criteria2.In, id)
}
func WhereCompany(id string) criteria2.ICondition {
	if id == "" {
		return nil
	}
	return criteria2.NewCondition("company", criteria2.Eq, id)
}
func WhereRole(id string) criteria2.ICondition {
	if id == "" {
		return nil
	}
	return criteria2.NewCondition("role", criteria2.Eq, id)
}
func WhereJobType(id string) criteria2.ICondition {
	if id == "" {
		return nil
	}
	return criteria2.NewCondition("job_type", criteria2.Eq, id)
}
