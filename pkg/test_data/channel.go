package test_data

import (
	"fmt"
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"strconv"
	"syreclabs.com/go/faker"
	"time"
)

func randomMinute() string {
	return faker.RandomChoice([]string{"*", "0", "1/5", "5,10,15"})
}
func randomHour() string {
	return faker.RandomChoice([]string{"*", "0", "1/5", "5,10,15,23"})
}
func randomDate() string {
	return faker.RandomChoice([]string{"*", "1", "5,10,15"})
}
func randomMonth() string {
	return strconv.Itoa(faker.RandomInt(1, 12))
}
func randomWeekday() string {
	return strconv.Itoa(faker.RandomInt(0, 6))
}
func ValidChannel(shouldActive bool) *channel.Entity {
	var usr user.Entity
	usr.SetID(uuid.New())
	title := faker.Company().Name()
	desc := faker.Lorem().Paragraph(5)
	active := shouldActive
	schedule := fmt.Sprintf("%s %s %s %s %s", randomMinute(), randomHour(), randomDate(), randomMonth(), randomWeekday())
	channels := []string{"email", "slack"}
	ch, _ := channel.New(
		uuid.New(),
		&usr,
		title,
		desc,
		active,
		schedule,
		time.Now(),
		time.Now(),
		channels,
		[]byte(`{"skills":"devops,aws,gcp", "keyword":"senior"}`),
	)
	return ch
}
