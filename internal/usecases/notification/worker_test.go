package notification

import (
	"testing"
)

/*
this will be under a loop. the responsibilities would be
queries unsent notifications. unsent notifications are notifications that are not sent, and is not part of schedule of the current time
notifications type are email and in-app (for now).
it should send a notifications
*/

func TestNotification_Worker(t *testing.T) {
	//strings.Title()
}
