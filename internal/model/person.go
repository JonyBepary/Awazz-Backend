package model

import (
	"fmt"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

var column []string = []string{
	"Id",
	"Attachment",
	"AttributedTo",
	"Context",
	"MediaType",
	"EndTime",
	"Generator",
	"Icon",
	"Image",
	"InReplyTo",
	"Location",
	"Preview",
	"Published",
	"Replies",
	"StartTime",
	"Summary",
	"Tag",
	"Updated",
	"Url",
	"Likes",
	"Shares",
	"Inbox",
	"Outbox",
	"Following",
	"Followers",
	"Liked",
	"PreferredUsername",
	"Endpoints",
	"Streams",
	"PublicKey"}
var column_type []string = []string{
	"INTEGER PRIMARY KEY",
	"BLOB",
	"TEXT",
	"TEXT",
	"TEXT",
	"NUMERIC",
	"TEXT",
	"BLOB",
	"BLOB",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"NUMERIC",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"TEXT",
	"INTEGER",
	"TEXT",
	"BLOB",
	"BLOB"}

func (p *Person) SavePerson() error {
	db, err := durable.CreateDatabase("Database/comments")
	if err != nil {
		return err
	}
	defer db.Close()
	fmt.Printf("CREATE TABLE IF NOT EXISTS PERSON ( %v)\n", durable.SetColumn(column, column_type))
	// db.Query(fmt.Sprintf("CREATE TABLE PERSON ( %v)", durable.SetColumn(column, column_type)))
	return nil
}
