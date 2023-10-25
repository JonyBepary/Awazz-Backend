package model

import (
	"fmt"
	"log"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (n *Notifications) SaveNotifications() error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")

	if err != nil {
		return err
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS NOTIFICATIONS (
	Receiver VARCHAR(128) PRIMARY KEY NOT NULL,
	Title VARCHAR(128),
    Body VARCHAR(128),
    Source VARCHAR(128),
    Image VARCHAR(128),
    Sound VARCHAR(128),
    Time INTEGER,
    Channel VARCHAR(128),
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
