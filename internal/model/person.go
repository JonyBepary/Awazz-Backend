package model

import (
	"fmt"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
	"github.com/davecgh/go-spew/spew"
)

func (p *Person) SavePerson(frag_num int64) error {

	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		return err
	}
	defer db.Close()

	sql_cmd := `CREATE TABLE IF NOT EXISTS PERSON (
	Id VARCHAR(255) PRIMARY KEY,
	Location VARCHAR(255),
	Attachment VARCHAR(255),
	AttributedTo VARCHAR(255),
	Context VARCHAR(255),
	MediaType VARCHAR(255),
	EndTime INTEGER,
	Generator VARCHAR(255),
	Icon VARCHAR(255),
	Image VARCHAR(255),
	InReplyTo VARCHAR(255),
	Preview VARCHAR(255),
	PublishedTime INTEGER,
	StartTime INTEGER,
	Summary VARCHAR(255),
	UpdatedTime INTEGER,
	Likes VARCHAR(255),
	Shares VARCHAR(255),
	Inbox VARCHAR(255),
	Outbox VARCHAR(255),
	PreferredUsername VARCHAR(255),
	PublicKey VARCHAR(255),
	FragmentationKey VARCHAR(255),
	Username VARCHAR(255),
	FOREIGN KEY (Username) REFERENCES USER(Username)
)`
	_, err = db.Exec(sql_cmd)
	if err != nil {
		return err
	}
	statement, err := db.Prepare("INSERT INTO PERSON (Id,Attachment,AttributedTo,Context,MediaType,EndTime,Generator,Icon,Image,InReplyTo,Location,Preview,PublishedTime,StartTime,Summary,UpdatedTime,Likes,Shares,Inbox,Outbox,PreferredUsername,PublicKey,FragmentationKey,Username) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(p.Id, p.Attachment, p.AttributedTo, p.Context, p.MediaType, p.EndTime, p.Generator, p.Icon, p.Image, p.InReplyTo, p.Location, p.Preview, p.PublishedTime, p.StartTime, p.Summary, p.UpdatedTime, p.Likes, p.Shares, p.Inbox, p.Outbox, p.PreferredUsername, p.PublicKey, p.FragmentationKey, p.Username)
	if err != nil {
		return err
	}
	return nil
}

func (p *Person) GetPerson(pid string, frag_num int64) error {

	// save fragmentation to file
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// spew.Dump(rows)
	//! fmt.Println("message id is: ", pid)
	row, err := db.Query("SELECT * FROM PERSON WHERE Id=?", pid)
	if err != nil {
		panic(err)
	}
	row.Next()
	err = row.Scan(&p.Id, &p.Attachment, &p.AttributedTo, &p.Context, &p.MediaType, &p.EndTime, &p.Generator, &p.Icon, &p.Image, &p.InReplyTo, &p.Location, &p.Preview, &p.PublishedTime, &p.Replies, &p.StartTime, &p.Summary, &p.UpdatedTime, &p.Url, &p.Too, &p.Bto, &p.Cc, &p.Bcc, &p.Likes, &p.Shares, &p.Inbox, &p.Outbox, &p.Following, &p.Followers, &p.Liked, &p.PreferredUsername, &p.Endpoints, &p.Streams, &p.PublicKey, &p.FragmentationKey, &p.Username)
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

func (p *Person) GetPersonByUsername(username string, frag_num int64) error {
	//! need to be fixed fragmentation
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// spew.Dump(rows)
	//! fmt.Println("message id is: ", pid)
	row, err := db.Query("SELECT * FROM PERSON WHERE Username=?", username)
	if err != nil {
		panic(err)
	}
	row.Next()
	err = row.Scan(&p.Id, &p.Attachment, &p.AttributedTo, &p.Context, &p.MediaType, &p.EndTime, &p.Generator, &p.Icon, &p.Image, &p.InReplyTo, &p.Location, &p.Preview, &p.PublishedTime, &p.Replies, &p.StartTime, &p.Summary, &p.UpdatedTime, &p.Url, &p.Too, &p.Bto, &p.Cc, &p.Bcc, &p.Likes, &p.Shares, &p.Inbox, &p.Outbox, &p.Following, &p.Followers, &p.Liked, &p.PreferredUsername, &p.Endpoints, &p.Streams, &p.PublicKey, &p.FragmentationKey, &p.Username)
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

func (p *Person) UpdatePerson(frag_num int64) error {

	// save fragmentation to file
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}

	defer db.Close()
	sql_cmd := `UPDATE PERSON SET Id=?,Attachment=?,AttributedTo=?,Context=?,MediaType=?,EndTime=?,Generator=?,Icon=?,Image=?,InReplyTo=?,Location=?,Preview=?,PublishedTime=?,StartTime=?,Summary=?,UpdatedTime=?,Likes=?,Shares=?,Inbox=?,Outbox=?,PreferredUsername=?,PublicKey=?,FragmentationKey=?,Username=? WHERE Id=?`
	statement, err := db.Prepare(sql_cmd)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Id, p.Attachment, p.AttributedTo, p.Context, p.MediaType, p.EndTime, p.Generator, p.Icon, p.Image, p.InReplyTo, p.Location, p.Preview, p.PublishedTime, p.StartTime, p.Summary, p.UpdatedTime, p.Likes, p.Shares, p.Inbox, p.Outbox, p.PreferredUsername, p.PublicKey, p.FragmentationKey, p.Username, p.Id)
	if err != nil {
		panic(err)
	}
	return nil
}

func (p *Person) DeletePerson(pid string, frag_num int64) error {

	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	spew.Dump(pid)
	//prepare statement
	statement, err := db.Prepare("DELETE FROM PERSON WHERE Id=?")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(pid)
	if err != nil {
		panic(err)
	}

	return nil
}

func (p *Person) GetPersonByFragmentationKey(fragmentationKey string, frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// spew.Dump(rows)
	//! fmt.Println("message id is: ", pid)
	row, err := db.Query("SELECT * FROM PERSON WHERE FragmentationKey=?", fragmentationKey)
	if err != nil {
		panic(err)
	}
	row.Next()
	err = row.Scan(&p.Id, &p.Attachment, &p.AttributedTo, &p.Context, &p.MediaType, &p.EndTime, &p.Generator, &p.Icon, &p.Image, &p.InReplyTo, &p.Location, &p.Preview, &p.PublishedTime, &p.Replies, &p.StartTime, &p.Summary, &p.UpdatedTime, &p.Url, &p.Too, &p.Bto, &p.Cc, &p.Bcc, &p.Likes, &p.Shares, &p.Inbox, &p.Outbox, &p.Following, &p.Followers, &p.Liked, &p.PreferredUsername, &p.Endpoints, &p.Streams, &p.PublicKey, &p.FragmentationKey, &p.Username)
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

// horizontally fragment this table by location
func (p *Person) FragmentateByLocation(frag_num int64) ([]*Person, error) {

	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM PERSON WHERE Location=?", p.Location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []*Person
	for rows.Next() {
		var person Person
		err = rows.Scan(&person.Id, &person.Attachment, &person.AttributedTo, &person.Context, &person.MediaType, &person.EndTime, &person.Generator, &person.Icon, &person.Image, &person.InReplyTo, &person.Location, &person.Preview, &person.PublishedTime, &person.Replies, &person.StartTime, &person.Summary, &person.UpdatedTime, &person.Url, &person.Too, &person.Bto, &person.Cc, &person.Bcc, &person.Likes, &person.Shares, &person.Inbox, &person.Outbox, &person.Following, &person.Followers, &person.Liked, &person.PreferredUsername, &person.Endpoints, &person.Streams, &person.PublicKey, &person.FragmentationKey, &person.Username)
		if err != nil {
			return nil, err
		}
		persons = append(persons, &person)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return persons, nil
}
