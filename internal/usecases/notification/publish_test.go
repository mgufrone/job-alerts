package notification

import "testing"

func TestNotification_Publish(t *testing.T) {
	// the method will always be called each time a new job entry is entered to database.
	// this method will first query channel which has the matched criteria related to it
	// if no criteria found, done
	// it would be considered a match if
	// all criterias under the channel should match with what the job have: skills, title, description, salary range, location, isRemote, etc
	// since channels are owned by a user, before creating a notification, it will have to check if the notification about the said job is already exists.
	// if it's not, create new notification that will be targeted to the user with the condition set in the channel (scheduled alert, buffered alert, immediate alert (disable this for now))
}
