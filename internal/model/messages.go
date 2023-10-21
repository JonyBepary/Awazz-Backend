package model

import (
	"fmt"
	"log"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (msg *Messages) SaveMessages() error {
	db, err := durable.CreateDatabase("Database/messages.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS Messages (
		MsgId VARCHAR(128) PRIMARY KEY,
		SenderId  VARCHAR(128),
		ReceiverId  VARCHAR(128),
		Content TEXT,
		SentTime TEXT,
		LastEdit TEXT,
		DeleteTime TEXT,
		Status BOOL,
		Attachment BLOB,
		Type TEXT,
		Reaction TEXT)
		`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}

	statement, err := db.Prepare("INSERT INTO Messages (MsgId,SenderId,ReceiverId,Content,SentTime,LastEdit,DeleteTime,Status,Type,Reaction) VALUES (?,?,?,?,?,?,?,?,?,?)")

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec(msg.MsgId, msg.SenderId, msg.ReceiverId, msg.Content, msg.SentTime.String(), msg.LastEdit.String(), msg.DeleteTime.String(), msg.Status, msg.Types, msg.Reaction)
	if err != nil {
		panic(err)
	}
	return nil
}

func (m *Messages) GetMessages(msgId string) error {
	db, err := durable.CreateDatabase("Database/messages.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select * from Messages where msgId = %v", msgId))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		SentTime := ""
		LastEdit := ""
		DeleteTime := ""
		err = rows.Scan(&m.MsgId, &m.SenderId, &m.ReceiverId, &m.Content, &SentTime, &LastEdit, &DeleteTime, &m.Status, &m.Attachment, &m.Types, &m.Reaction)

		if err != nil {
			log.Fatal(err)
		}

		t, err := time.Parse(time.RFC3339, SentTime)
		if err != nil {
			panic(err)
		}

		m.SentTime = timestamppb.New(t)

		t, err = time.Parse(time.RFC3339, LastEdit)
		if err != nil {
			panic(err)
		}

		m.LastEdit = timestamppb.New(t)

		t, err = time.Parse(time.RFC3339, DeleteTime)
		if err != nil {
			panic(err)
		}

		m.DeleteTime = timestamppb.New(t)

	}
	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (u *Messages) UpdatedMessages(msgId string) error {
	db, err := durable.CreateDatabase("Database/messages.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	u.LastEdit = timestamppb.Now()
	_, err = db.Exec("UPDATE Messages SET  LastEdit= ?, Content = ? WHERE 	msId = ? ", u.LastEdit, u.Content, msgId)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
