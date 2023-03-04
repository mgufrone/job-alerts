package test_data

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/internal/domain/user_channel"
	"syreclabs.com/go/faker"
	"time"
)

var (
	receivers = map[string]string{
		"email":    "test@test.com",
		"telegram": "1232323",
		"slack":    "https://webhook.slack.com/somewherethere",
	}
)

func ValidUserChannel(usr *user.Entity, chType string) *user_channel.Entity {
	if chType == "" {
		chType = faker.RandomChoice([]string{"email", "telegram", "slack"})
	}
	chTarget := receivers[chType]
	ch, _ := user_channel.New(
		uuid.New(),
		usr,
		chType,
		chTarget,
		true,
		time.Now(),
		time.Now(),
	)
	return ch
}
