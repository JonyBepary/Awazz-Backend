package model

import (
	"fmt"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
	"github.com/SohelAhmedJoni/Awazz-Backend/pkg"
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
	"TEXT",
	"TEXT",
	"BLOB",
	"BLOB",
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
	db, err := durable.CreateDatabase("./Database/persons.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS PERSON ( %v)", durable.SetColumn(column, column_type)))
	if err != nil {
		return err
	}
	statement, err := db.Prepare("INSERT INTO PERSON (Id,AttributedTo,Context,MediaType,EndTime,Generator,InReplyTo,Location,Preview,Published,Replies,StartTime,Summary,Tag,Updated,Url,Likes,Shares,Inbox,Outbox,Following,Followers,Liked,PreferredUsername,Endpoints) VALUES (? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? )")
	if err != nil {
		return err
	}
	_, err = statement.Exec(p.Id, pkg.ReadFile(p.Attachment), p.AttributedTo, p.Context, p.MediaType, p.EndTime.String(), p.Generator, pkg.ReadFile(p.Icon), pkg.ReadFile(p.Image), p.InReplyTo, p.Location, p.Preview, p.Published.String(), p.Replies, p.StartTime.String(), p.Summary, p.Tag[0], p.Updated.String(), p.Url, p.To[0], p.Bto[0], p.Cc[0], p.Bcc[0], p.Likes, p.Shares, p.Inbox, p.Outbox, p.Following, p.Followers, p.Liked, p.PreferredUsername, p.Endpoints, p.Streams[0], p.PublicKey)
	if err != nil {
		return err
	}
	return nil
}
