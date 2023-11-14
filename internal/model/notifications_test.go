package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotifications_SaveNotifications(t *testing.T) {
	// Create a new Notifications object with some data
	n := &Notifications{
		Receiver:      "jony",
		Title:         "Test Notification",
		Body:          "This is a test notification",
		Source:        "test",
		Image:         "test.png",
		Sound:         "test.mp3",
		Time:          1234567890,
		Channel:       "test-channel",
		PriorityLevel: 1,
		ReadStatus:    false,
		Created:       1234567890,
	}

	// Call the SaveNotifications method
	err := n.SaveNotifications()

	// Check that there were no errors
	assert.NoError(t, err)
}
