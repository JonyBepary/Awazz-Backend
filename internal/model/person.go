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
	"PublicKey",
	"FragmentationKey"}
var column_type []string = []string{
	"VARCHAR(256) PRIMARY KEY",
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
	"TEXT"}

func (p *Person) SavePerson() error {
	db, err := durable.CreateDatabase("./Database/", p.FragmentationKey, "person.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS PERSON ( %v)", durable.SetColumn(column, column_type)))
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare("INSERT INTO PERSON (Id,AttributedTo,Context,MediaType,EndTime,Generator,InReplyTo,Location,Preview,Published,Replies,StartTime,Summary,Tag,Updated,Url,Likes,Shares,Inbox,Outbox,Following,Followers,Liked,PreferredUsername,Endpoints,FragmentationKey) VALUES (? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,?,? )")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Id, p.AttributedTo, p.Context, p.MediaType, p.EndTime.String(), p.Generator, p.Icon, p.Image, p.InReplyTo, p.Location, p.Preview, p.Published.String(), p.Replies, p.StartTime.String(), p.Summary, nil, p.Updated.String(), p.Url, p.Likes, p.Shares, p.Inbox, p.Outbox, p.Following, p.Followers, p.Liked, p.PreferredUsername, p.Endpoints, p.FragmentationKey)
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *Person) GetPerson(pid string) error {
	db,
		err := durable.CreateDatabase("Database/person")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return nil
}
func (p *Person) UpdatePerson(pid string) error {
	db, err := durable.CreateDatabase("Database/person")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return nil
}
func (p *Person) DeletePerson(pid string) error {
	db, err := durable.CreateDatabase("Database/person")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return nil
}
