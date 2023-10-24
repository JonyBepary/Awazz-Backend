package model

import (
	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (p *Person) SavePerson() error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()
	sql_cmd := `CREATE TABLE IF NOT EXISTS PERSON (
	Id VARCHAR(255) PRIMARY KEY,
	Attachment VARCHAR(255),
	AttributedTo VARCHAR(255),
	Context VARCHAR(255),
	MediaType VARCHAR(255),
	EndTime INTEGER,
	Generator VARCHAR(255),
	Icon VARCHAR(255),
	Image VARCHAR(255),
	InReplyTo VARCHAR(255),
	Location VARCHAR(255),
	Preview VARCHAR(255),
	PublishedTime INTEGER,
	Replies VARCHAR(255),
	StartTime INTEGER,
	Summary VARCHAR(255),
	Tag VARCHAR(255),
	UpdatedTime INTEGER,
	Url VARCHAR(255),
	Too VARCHAR(255),
	Bto VARCHAR(255),
	Cc VARCHAR(255),
	Bcc VARCHAR(255),
	Likes VARCHAR(255),
	Shares VARCHAR(255),
	Inbox VARCHAR(255),
	Outbox VARCHAR(255),
	Following VARCHAR(255),
	Followers VARCHAR(255),
	Liked VARCHAR(255),
	PreferredUsername VARCHAR(255),
	Endpoints VARCHAR(255),
	Streams VARCHAR(255),
	PublicKey VARCHAR(255),
	FragmentationKey VARCHAR(255)
)`
	_, err = db.Exec(sql_cmd)
	if err != nil {
		return err
	}
	statement, err := db.Prepare("INSERT INTO PERSON (Id,Attachment,AttributedTo,Context,MediaType,EndTime,Generator,Icon,Image,InReplyTo,Location,Preview,PublishedTime,Replies,StartTime,Summary,Tag,UpdatedTime,Url,Too,Bto,Cc,Bcc,Likes,Shares,Inbox,Outbox,Following,Followers,Liked,PreferredUsername,Endpoints,Streams,PublicKey,FragmentationKey) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(p.Id, p.Attachment, p.AttributedTo, p.Context, p.MediaType, p.EndTime, p.Generator, p.Icon, p.Image, p.InReplyTo, p.Location, p.Preview, p.PublishedTime, p.Replies, p.StartTime, p.Summary, p.Tag, p.UpdatedTime, p.Url, p.Too, p.Bto, p.Cc, p.Bcc, p.Likes, p.Shares, p.Inbox, p.Outbox, p.Following, p.Followers, p.Liked, p.PreferredUsername, p.Endpoints, p.Streams, p.PublicKey, p.FragmentationKey)
	if err != nil {
		return err
	}
	return nil
}

func (p *Person) GetPerson(msgId string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// spew.Dump(rows)
	//! fmt.Println("message id is: ", msgId)
	row, err := db.Query("SELECT * FROM PERSON WHERE Id=?", msgId)
	if err != nil {
		panic(err)
	}
	row.Next()
	err = row.Scan(&p.Id, &p.Attachment, &p.AttributedTo, &p.Context, &p.MediaType, &p.EndTime, &p.Generator, &p.Icon, &p.Image, &p.InReplyTo, &p.Location, &p.Preview, &p.PublishedTime, &p.Replies, &p.StartTime, &p.Summary, &p.Tag, &p.UpdatedTime, &p.Url, &p.Too, &p.Bto, &p.Cc, &p.Bcc, &p.Likes, &p.Shares, &p.Inbox, &p.Outbox, &p.Following, &p.Followers, &p.Liked, &p.PreferredUsername, &p.Endpoints, &p.Streams, &p.PublicKey, &p.FragmentationKey)
	if err != nil {
		panic(err)
	}

	err = row.Err()
	if err != nil {
		panic(err)
	}
	row.Close()

	//! spew.Dump(p.Id)
	return nil
}

// func GetPerson(msgId string) (Person, error) {
// 	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	// spew.Dump(rows)
// 	//! fmt.Println("message id is: ", msgId)
// 	row, err := db.Query("SELECT * FROM PERSON WHERE Id=?", msgId)
// 	if err != nil {
// 		panic(err)
// 	}
// 	p := Person{}
// 	row.Next()
// 	err = row.Scan(&p.Id, &p.Attachment, &p.AttributedTo, &p.Context, &p.MediaType, &p.EndTime, &p.Generator, &p.Icon, &p.Image, &p.InReplyTo, &p.Location, &p.Preview, &p.PublishedTime, &p.Replies, &p.StartTime, &p.Summary, &p.Tag, &p.UpdatedTime, &p.Url, &p.Too, &p.Bto, &p.Cc, &p.Bcc, &p.Likes, &p.Shares, &p.Inbox, &p.Outbox, &p.Following, &p.Followers, &p.Liked, &p.PreferredUsername, &p.Endpoints, &p.Streams, &p.PublicKey, &p.FragmentationKey)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = row.Err()
// 	if err != nil {
// 		panic(err)
// 	}
// 	row.Close()

// 	//! spew.Dump(p.Id)
// 	return p, nil
// }
