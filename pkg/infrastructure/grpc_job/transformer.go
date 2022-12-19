package grpc_job

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"time"
)

func Transform(ent *job.Entity) *Job {
	var (
		salaries []float32
	)
	for _, sal := range ent.Salary() {
		salaries = append(salaries, float32(sal))
	}
	return &Job{
		Id:             ent.ID().String(),
		Role:           ent.Role(),
		CompanyName:    ent.CompanyName(),
		CompanyUrl:     ent.CompanyURL(),
		Description:    ent.Description(),
		JobUrl:         ent.JobURL(),
		Tags:           ent.Tags(),
		Location:       ent.Location(),
		JobType:        ent.JobType(),
		IsRemote:       ent.IsRemote(),
		Salary:         salaries,
		SalaryCurrency: ent.SalaryCurrency(),
		CreatedAt:      ent.CreatedAt().Unix(),
		UpdatedAt:      ent.UpdatedAt().Unix(),
		Source:         ent.Source(),
	}
}
func ToJob(prop *Job) (*job.Entity, error) {
	id := prop.GetId()
	if id == "" {
		id = uuid.Nil.String()
	}
	j, err := job.NewJob(
		uuid.MustParse(id),
		prop.GetRole(),
		prop.GetCompanyName(),
		prop.GetCompanyUrl(),
		prop.GetDescription(),
		prop.GetJobUrl(),
		prop.GetSource(),
		prop.GetLocation(),
		prop.GetTags(),
	)
	if err != nil {
		return nil, err
	}
	j.SetJobType(prop.GetJobType())
	j.SetIsRemote(prop.GetIsRemote())
	propSalary := prop.GetSalary()
	var salary []float64
	for _, v := range propSalary {
		salary = append(salary, float64(v))
	}
	j.SetSalary(salary)
	j.SetSalaryCurrency(prop.GetSalaryCurrency())
	j.SetUpdatedAt(time.Unix(prop.GetUpdatedAt(), 0))
	j.SetCreatedAt(time.Unix(prop.GetCreatedAt(), 0))
	return j, err
}
func ToSliceJob(props []*Job) (out []*job.Entity, err error) {
	for _, r := range props {
		pr, err1 := ToJob(r)
		if err1 != nil {
			err = err1
			return
		}
		out = append(out, pr)
	}
	return
}
