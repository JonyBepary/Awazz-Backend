package model
func TestSaveMessages(t *testing.T) {
	msg := &Messages{
		MsgId:      "123",
		SenderId:   "sender",
		ReceiverId: "receiver",
		Types:      "text",
		Content:    "hello",
		SentTime:   123456,
		LastEdit:   123456,
		DeleteTime: "2022-01-01",
		Status:     true,
		Attachment: "attachment",
		Reaction:   "reaction",
	}

	err := msg.SaveMessages(1)
	if err != nil {
		t.Errorf("SaveMessages returned an error: %v", err)
	}

	// Verify that the message was saved correctly
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_1.sqlite")
	if err != nil {
		t.Errorf("Failed to open database: %v", err)
	}
	defer db.Close()

	var savedMsg Messages
	err = db.QueryRow("SELECT * FROM MESSAGE WHERE MsgId=?", msg.MsgId).Scan(
		&savedMsg.MsgId,
		&savedMsg.SenderId,
		&savedMsg.ReceiverId,
		&savedMsg.Types,
		&savedMsg.Content,
		&savedMsg.SentTime,
		&savedMsg.LastEdit,
		&savedMsg.DeleteTime,
		&savedMsg.Status,
		&savedMsg.Attachment,
		&savedMsg.Reaction,
	)
	if err != nil {
		t.Errorf("Failed to retrieve saved message: %v", err)
	}

	if !reflect.DeepEqual(msg, &savedMsg) {
		t.Errorf("Saved message does not match original message. Expected: %v, got: %v", msg, &savedMsg)
	}
}
