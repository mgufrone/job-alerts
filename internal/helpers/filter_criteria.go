package helpers

import (
	"encoding/json"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"strings"
)

func FilterToCriteria(ch *channel.Entity, cb criteria.ICriteriaBuilder) criteria.ICriteriaBuilder {
	var (
		cr     map[string]interface{}
		wheres []criteria.ICriteriaBuilder
	)
	json.Unmarshal(ch.Criterias(), &cr)
	for k, v := range cr {
		switch k {
		case "skills":
			var val []string
			switch vt := v.(type) {
			case string:
				val = strings.Split(vt, ",")
			case []string:
				val = vt
			}
			wheres = append(wheres, cb.Where(job.WhereTags(val)))
		case "keyword":
			switch vt := v.(type) {
			case string:
				wheres = append(wheres, cb.Or(
					cb.Where(job.WhereRoleContains(vt)),
					cb.Where(job.WhereDescriptionContains(vt)),
				))
			}
		}
	}
	return cb.And(wheres...)
}
