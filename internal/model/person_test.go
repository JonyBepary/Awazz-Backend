package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSavePerson(t *testing.T) {
	// Create a new Person object with test data
	p := &Person{
		Id:                "test_id",
		Location:          "test_location",
		Attachment:        "test_attachment",
		AttributedTo:      "test_attributed_to",
		Context:           "test_context",
		MediaType:         "test_media_type",
		EndTime:           1234567890,
		Generator:         "test_generator",
		Icon:              "test_icon",
		Image:             "test_image",
		InReplyTo:         "test_in_reply_to",
		Preview:           "test_preview",
		PublishedTime:     1234567890,
		StartTime:         1234567890,
		Summary:           "test_summary",
		UpdatedTime:       1234567890,
		Likes:             "test_likes",
		Shares:            "test_shares",
		Inbox:             "test_inbox",
		Outbox:            "test_outbox",
		PreferredUsername: "test_preferred_username",
		PublicKey:         "test_public_key",
		FragmentationKey:  "test_fragmentation_key",
		Username:          "test_username",
	}

	// Call the SavePerson method with the test data
	err := p.SavePerson(1)

	// Check that the method returned no errors
	assert.NoError(t, err)
}
