package model

import (
	"fmt"
	"log"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (msg *Messages) SaveMessages(frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS MESSAGE (
		MsgId VARCHAR(255) PRIMARY KEY,
		SenderId  VARCHAR(255),
		ReceiverId  VARCHAR(255),
		Types TEXT,
		Content TEXT,
		SentTime INTEGER,
		LastEdit INTEGER,
		DeleteTime TEXT,
		Status BOOL,
		Attachment TEXT,
		Reaction TEXT)
	`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}

	statement, err := db.Prepare("INSERT INTO MESSAGE (MsgId,SenderId,ReceiverId,Content,SentTime,LastEdit,DeleteTime,Status,Attachment,Types,Reaction) VALUES (?,?,?,?,?,?,?,?,?,?,?)")

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec(msg.MsgId, msg.SenderId, msg.ReceiverId, msg.Content, msg.SentTime, msg.LastEdit, msg.DeleteTime, msg.Status, msg.Attachment, msg.Types, msg.Reaction)
	if err != nil {
		panic(err)
	}
	return nil
}

func (m *Messages) GetMessages(msgId string, frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select * from MESSAGE where msgId = %v", msgId))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&m.MsgId, &m.SenderId, &m.ReceiverId, &m.Content, &m.SentTime, &m.LastEdit, &m.DeleteTime, &m.Status, &m.Attachment, &m.Types, &m.Reaction)

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

func (u *Messages) UpdatedMessages(msgId string, frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	u.LastEdit = time.Now().Unix()
	_, err = db.Exec("UPDATE Messages SET  LastEdit= ?, Content = ? WHERE 	msId = ? ", u.LastEdit, u.Content, msgId)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (d *Messages) DeleteMessages(msgId string, frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM  MESSAGE WHERE  msgId= ?", msgId)

	d.LastEdit = time.Now().Unix()

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
