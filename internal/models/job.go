package models

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/db"
)

type Tag struct {
	db.Entity
	Name string `gorm:"size:100;uniqueIndex"`
	Jobs []*Job `gorm:"many2many:job_tags"`
}
type JobTag struct {
	JobID uuid.UUID `gorm:"primaryKey"`
	TagID uuid.UUID `gorm:"primaryKey"`
	db.EntityTimestamp
}

type Job struct {
	db.Entity
	Role        string `gorm:"index"`
	CompanyName string `gorm:"index"`
	CompanyURL  string `gorm:"index"`
	Description string `gorm:"type:text"`
	JobURL      string
	Tags        []*Tag `gorm:"many2many:job_tags;"`
	Location    string `gorm:"index"`
	Source      string `gorm:"index"`
	// can be optional
	JobType        string  `gorm:"index"`
	IsRemote       bool    `gorm:"index"`
	SalaryMin      float64 `gorm:"index"`
	SalaryMax      float64 `gorm:"index"`
	SalaryCurrency string  `gorm:"index"`
	// timestamps
	db.EntityTimestamp
}

func (j *Job) FromDomain(job2 *job.Entity) {
	salaries := job2.Salary()
	j.SalaryMin = 0
	if len(salaries) != 0 {
		j.SalaryMin = salaries[0]
	}
	if len(salaries) > 1 {
		j.SalaryMax = salaries[len(salaries)-1]
	}
	tags := make([]*Tag, len(job2.Tags()))
	for idx, t := range job2.Tags() {
		tags[idx] = &Tag{Name: t}
	}
	j.ID = job2.ID()
	j.Role = job2.Role()
	j.CompanyName = job2.CompanyName()
	j.CompanyURL = job2.CompanyURL()
	j.Description = job2.Description()
	j.Source = job2.Source()
	j.Location = job2.Location()
	j.Tags = tags
	j.JobURL = job2.JobURL()
	j.JobType = job2.JobType()
	j.SalaryCurrency = job2.SalaryCurrency()
	j.CreatedAt = job2.CreatedAt()
	j.UpdatedAt = job2.UpdatedAt()
	j.IsRemote = job2.IsRemote()
}
func (j *Job) Transform() (*job.Entity, error) {
	if j == nil {
		return nil, nil
	}
	tags := make([]string, len(j.Tags))
	for idx, t := range j.Tags {
		tags[idx] = t.Name
	}
	jb, err := job.NewJob(
		j.ID,
		j.Role,
		j.CompanyName,
		j.CompanyURL,
		j.Description,
		j.JobURL,
		j.Source,
		j.Location,
		tags,
	)
	if err != nil {
		return nil, err
	}
	jb.SetJobType(j.JobType)
	salaries := []float64{j.SalaryMin}
	if j.SalaryMax > 0 {
		salaries = append(salaries, j.SalaryMax)
	}
	jb.SetSalary(salaries)
	jb.SetSalaryCurrency(j.SalaryCurrency)
	jb.SetCreatedAt(j.CreatedAt)
	jb.SetUpdatedAt(j.UpdatedAt)
	jb.SetIsRemote(j.IsRemote)
	return jb, err
}
