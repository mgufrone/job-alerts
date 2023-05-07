//go:generate go run mgufrone.dev/job-alerts/cmd/generate-domain mgufrone.dev/job-alerts/internal/domain/job mgufrone.dev/job-alerts
package job

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/mgufrone/go-utils/try"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	errors3 "mgufrone.dev/job-alerts/pkg/errors"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Entity struct {
	id          uuid.UUID
	role        string
	companyName string
	companyURL  string
	description string
	jobURL      string
	tags        []string
	location    string
	source      string
	// can be optional
	jobType        string
	isRemote       bool
	salary         []float64
	salaryCurrency string
	// timestamps
	createdAt time.Time
	updatedAt time.Time
	title     string
}

func (j *Entity) Id() uuid.UUID {
	return j.id
}

func (j *Entity) SetId(id uuid.UUID) {
	j.id = id
}

func (j *Entity) Title() string {
	return j.title
}

func (j *Entity) SetTitle(title string) (err error) {
	j.title = title
	return nil
}

func (j *Entity) Source() string {
	return j.source
}

func (j *Entity) SetSource(source string) (err error) {
	if err = validation.Validate(source,
		validation.Required, validation.NotNil); err != nil {
		return errors3.FieldError("source", err)
	}
	j.source = source
	return
}

func NewJob(
	id uuid.UUID, role, companyName, companyURL, description, jobURL, source, location string,
	tags []string, title string) (*Entity, error) {
	job := &Entity{}
	err := try.RunOrError(func() error {
		return job.SetID(id)
	}, func() error {
		return job.SetRole(role)
	}, func() error {
		return job.SetCompanyName(companyName)
	}, func() error {
		return job.SetCompanyURL(companyURL)
	}, func() error {
		return job.SetDescription(description)
	}, func() error {
		return job.SetJobURL(jobURL)
	}, func() error {
		return job.SetTags(tags)
	}, func() error {
		return job.SetLocation(location)
	}, func() error {
		return job.SetSource(source)
	}, func() error {
		return job.SetTitle(title)
	})
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j *Entity) ID() uuid.UUID {
	return j.id
}

func (j *Entity) SetID(id uuid.UUID) (err error) {
	if err = try.RunOrError(func() error {
		return validation.Validate(id,
			validation.Required, validation.NotNil,
		)
	}); err != nil {
		return errors3.FieldError("id", err)
	}
	if id == uuid.Nil {
		err = errors3.FieldError("id", fmt.Errorf("invalid value"))
		return
	}
	j.id = id
	return
}

func (j *Entity) Role() string {
	return j.role
}

func (j *Entity) SetRole(role string) (err error) {
	regex := regexp.MustCompile(`(?i)[^\w\s-,]+`)
	space := regexp.MustCompile(`\s+`)
	role = regex.ReplaceAllString(role, " ")
	role = space.ReplaceAllString(role, " ")
	role = cases.Title(language.Und, cases.NoLower).String(role)
	role = strings.TrimSpace(role)
	if err = try.RunOrError(func() error {
		return validation.Validate(role,
			validation.Required, validation.NotNil,
		)
	}); err != nil {
		return errors3.FieldError("role", err)
	}
	j.role = role
	return
}

func (j *Entity) CompanyName() string {
	return j.companyName
}

func (j *Entity) SetCompanyName(companyName string) (err error) {
	if err = try.RunOrError(func() error {
		return validation.Validate(companyName,
			validation.Required, validation.NotNil,
		)
	}); err != nil {
		return errors3.FieldError("company_name", err)
	}
	j.companyName = companyName
	return
}

func (j *Entity) CompanyURL() string {
	return j.companyURL
}

func (j *Entity) SetCompanyURL(companyURL string) (err error) {
	if companyURL != "" && companyURL != "-" {
		if err = validation.Validate(companyURL, validation.Required, is.URL); err != nil {
			return errors3.FieldError("company_url", err)
		}
	}
	j.companyURL = companyURL
	return
}

func (j *Entity) Description() string {
	return j.description
}

func (j *Entity) SetDescription(description string) (err error) {
	if err = validation.Validate(description, validation.Required, validation.NotNil); err != nil {
		return errors3.FieldError("description", err)
	}
	j.description = description
	return
}

func (j *Entity) Tags() []string {
	return j.tags
}

func (j *Entity) SetTags(tags []string) (err error) {
	j.tags = tags
	return
}

func (j *Entity) Salary() []float64 {
	return j.salary
}

func (j *Entity) SetSalary(salary []float64) (err error) {
	j.salary = salary
	return
}

func (j *Entity) SalaryCurrency() string {
	return j.salaryCurrency
}

func (j *Entity) SetSalaryCurrency(salaryCurrency string) (err error) {
	j.salaryCurrency = salaryCurrency
	return
}

func (j *Entity) Location() string {
	return j.location
}

func (j *Entity) SetLocation(location string) (err error) {
	j.location = location
	return
}

func (j *Entity) JobURL() string {
	return j.jobURL
}

func (j *Entity) SetJobURL(jobURL string) (err error) {
	j.jobURL, err = url.QueryUnescape(jobURL)
	return
}

func (j *Entity) JobType() string {
	return j.jobType
}

func (j *Entity) SetJobType(jobType string) (err error) {
	j.jobType = jobType
	return
}

func (j *Entity) IsRemote() bool {
	return j.isRemote
}

func (j *Entity) SetIsRemote(isRemote bool) (err error) {
	j.isRemote = isRemote
	return
}

func (j *Entity) CreatedAt() time.Time {
	return j.createdAt
}

func (j *Entity) SetCreatedAt(createdAt time.Time) (err error) {
	j.createdAt = createdAt
	return
}

func (j *Entity) UpdatedAt() time.Time {
	return j.updatedAt
}

func (j *Entity) SetUpdatedAt(updatedAt time.Time) (err error) {
	j.updatedAt = updatedAt
	return
}
