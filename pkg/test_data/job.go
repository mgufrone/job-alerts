package test_data

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"strconv"
	"syreclabs.com/go/faker"
)

func ValidJob() *job.Entity {
	role := faker.RandomChoice([]string{"software engineer", "marketing manager", "devops engineer", "IT Manager", "Scrum Master"})
	jobSource := faker.RandomChoice([]string{"tia", "kalibrr", "hubstaff", "upwork", "offline", "flyer", "google", "jobstreet"})
	city := faker.RandomChoice([]string{"city a", "city b"})
	currency := faker.RandomChoice([]string{"idr", "sgd", "thb"})
	tags := faker.Lorem().Words(3)
	salaryMin, _ := strconv.ParseFloat(faker.Number().Decimal(10, 3), 0)
	salaryMax, _ := strconv.ParseFloat(faker.Number().Decimal(10, 3), 0)
	jb, _ := job.NewJob(
		uuid.New(),
		role,
		faker.Company().Name(),
		faker.Internet().Url(),
		faker.Lorem().Paragraph(10),
		faker.Internet().Url(),
		jobSource,
		city,
		tags,
	)
	jb.SetSalary([]float64{salaryMin, salaryMax + 1})
	jb.SetSalaryCurrency(currency)
	jb.SetJobType(faker.RandomChoice([]string{"full-time", "contract", "freelance"}))
	jb.SetIsRemote(faker.RandomInt(0, 1) == 1)
	return jb
}
