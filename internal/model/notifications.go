package model

import (
	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (n *Notifications) SaveNotifications() error {
	db, err := durable.CreateDatabase("Database/notifications")
	if err != nil {
		return err
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS Messages (
	Title VARCHAR(128),
    Body VARCHAR(128),
    Source VARCHAR(128),
    Image VARCHAR(128),
    Sound VARCHAR(128),
    Time TIMESTAMP,
    Channel VARCHAR(128),
    PriorityLevel INT,
    ReadStatus bool,
    Created TIMESTAMP)
	`
	_, err = db.Exec(str)
	if err != nil {
		return err
	}
	str2 := `
INSERT INTO Messages (Title,Body,Source,Image,Sound,Time,Channel,PriorityLevel,ReadStatus,Created) VALUES (?,?,?,?,?,?,?,?,?,?);
	`
	statement, err := db.Prepare(str2)
	if err != nil {
		return err
	}
	_, err = statement.Exec(n.Title, n.Body, n.Source, n.Image, n.Sound, n.Time, n.Channel, n.PriorityLevel, n.ReadStatus, n.Created)
	if err != nil {
		return err
	}

	return nil
}

/*func (n *Notifications) GetNotifications() error {
	db, err := durable.CreateDatabase("Database/notifications")
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}
func (n *Notifications) UpdateNotifications() error {
	db, err := durable.CreateDatabase("Database/notifications")
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}
func (n *Notifications) DeleteNotifications() error {
	db, err := durable.CreateDatabase("Database/notifications")
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}
*/
