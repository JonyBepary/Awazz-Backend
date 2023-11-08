package model

import (
	"fmt"
	"log"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

//create table
func (n *Notifications) SaveNotifications() error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")

	if err != nil {
		return err
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS NOTIFICATIONS (
	Receiver VARCHAR(255) PRIMARY KEY NOT NULL,
	Title VARCHAR(255),
    Body VARCHAR(255),
    Source VARCHAR(255),
    Image VARCHAR(255),
    Sound VARCHAR(255),
    Time INTEGER,
    Channel VARCHAR(255),
    PriorityLevel INTEGER,
    ReadStatus bool,
    Created INTEGER)
	`
	_, err = db.Exec(str)
	if err != nil {
		return err
	}
	str2 := `
INSERT INTO NOTIFICATIONS (Receiver,Title,Body,Source,Image,Sound,Time,Channel,PriorityLevel,ReadStatus,Created) VALUES (?,?,?,?,?,?,?,?,?,?,?);
	`
	statement, err := db.Prepare(str2)
	if err != nil {
		return err
	}
	_, err = statement.Exec(n.Receiver, n.Title, n.Body, n.Source, n.Image, n.Sound, n.Time, n.Channel, n.PriorityLevel, n.ReadStatus, n.Created)
	if err != nil {
		return err
	}

	return nil
}
/// read table
func (n *Notifications) GetNotifications(msgId string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select * from Notifications where msgId = %v", msgId))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&n.Title, &n.Body, &n.Source, &n.Image, &n.Sound, &n.Time, &n.Channel, &n.PriorityLevel, &n.ReadStatus, &n.Created)

		if err != nil {
			log.Fatal(err)
		}

	}
	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (u *Notifications) UpdateNotifications(Title, Body, Source string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	u.Time = time.Now().Unix()
	u.Title = Title
	u.Body = Body
	u.Source = Source
	_, err = db.Exec("UPDATE Notification SET  Time= ?, Title = ?, Body = ? WHERE Source = ? ", u.Time, u.Title, u.Body, u.Source)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (d *Notifications) DeleteNotifications(Source string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM  Notifications WHERE  Source= ?", Source)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
