package model

import "mgufrone.dev/job-alerts/internal/domain/job"

func (j *Job) FromDomain(ent *job.Entity) {
	j.ID = ent.ID().String()
	j.Role = ent.Role()
	j.Description = ent.Description()
	j.Tags = ent.Tags()
	j.Salary = ent.Salary()
	j.CreatedAt = ent.CreatedAt()
	j.UpdatedAt = ent.UpdatedAt()
	j.IsRemote = ent.IsRemote()
	j.CompanyURL = ent.CompanyURL()
	j.CompanyName = ent.CompanyName()
	j.URL = ent.JobURL()
}
