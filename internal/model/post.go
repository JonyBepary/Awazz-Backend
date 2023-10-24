package model

import (
	"log"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (p *Post) SavePost() error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE if not EXISTS POST(
    Id VARCHAR(255) PRIMARY KEY,
    Community VARCHAR(255),
    Content TEXT,
    CreatedAt INTEGER,
    UpdatedAt INTEGER,
    DeletedAt INTEGER,
    Likes INTEGER,
    Shares INTEGER,
    Comments INTEGER,
    Author VARCHAR(255),
    Parent VARCHAR(255),
    Rank INTEGER,
    Children VARCHAR(255),
    Tags VARCHAR(255),
    Mentions VARCHAR(255),
    IsSensitive BOOLEAN,
    IsNsfw BOOLEAN,
    IsDeleted BOOLEAN,
    IsPinned BOOLEAN,
    IsEdited BOOLEAN,
    IsLiked BOOLEAN,
    IsShared BOOLEAN,
    IsCommented BOOLEAN,
    IsSubscribed BOOLEAN,
    IsBookmarked BOOLEAN,
    IsReblogged BOOLEAN,
    IsMentioned BOOLEAN,
    IsPoll BOOLEAN,
    IsPollVoted BOOLEAN,
    IsPollExpired BOOLEAN,
    IsPollClosed BOOLEAN,
    IsPollMultiple BOOLEAN,
    IsPollHideTotals BOOLEAN,
    FragmentationKey VARCHAR(255)
)

	`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}
	str2 := `INSERT INTO Post (Id,Community,Content,CreatedAt,UpdatedAt,DeletedAt,Likes,Shares,Comments,Author,Parent,Rank,IsSensitive,IsNsfw,IsDeleted,IsPinned,IsEdited,IsLiked,IsShared,IsCommented,IsSubscribed,IsBookmarked,IsReblogged,IsMentioned,IsPoll,IsPollVoted,IsPollExpired,IsPollClosed,IsPollMultiple,IsPollHideTotals,FragmentationKey) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	statement, err := db.Prepare(str2)
	if err != nil {
		panic(err)
	}
	_, err = (statement.Exec(p.Id,
		p.Community,
		p.Content,
		p.CreatedAt,
		p.UpdatedAt,
		p.DeletedAt,
		p.Likes,
		p.Shares,
		p.Comments,
		p.Author,
		p.Parent,
		p.Rank,
		p.Children,
		p.Tags,
		p.Mentions,
		p.IsSensitive,
		p.IsNsfw,
		p.IsDeleted,
		p.IsPinned,
		p.IsEdited,
		p.IsLiked,
		p.IsShared,
		p.IsCommented,
		p.IsSubscribed,
		p.IsBookmarked,
		p.IsReblogged,
		p.IsMentioned,
		p.IsPoll,
		p.IsPollVoted,
		p.IsPollExpired,
		p.IsPollClosed,
		p.IsPollMultiple,
		p.IsPollHideTotals,
		p.FragmentationKey,
	))
	if err != nil {
		panic(err)
	}

	return nil
}
func (p *Post) GetPost(msgId string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// spew.Dump(rows)
	rows, err := db.Query("SELECT * FROM ")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(p.Id, p.Community, p.Content, p.CreatedAt, p.UpdatedAt, p.DeletedAt, p.Likes, p.Shares, p.Comments, p.Author, p.Parent, p.Rank, p.Children, p.Tags, p.Mentions, p.IsSensitive, p.IsNsfw, p.IsDeleted, p.IsPinned, p.IsEdited, p.IsLiked, p.IsShared, p.IsCommented, p.IsSubscribed, p.IsBookmarked, p.IsReblogged, p.IsMentioned, p.IsPoll, p.IsPollVoted, p.IsPollExpired, p.IsPollClosed, p.IsPollMultiple, p.IsPollHideTotals, p.FragmentationKey)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return nil
}
func (u *Post) UpdatePost(ID string) error {
	db, err := durable.CreateDatabase("./Database/", "common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	u.UpdatedAt = time.Now().Unix()

	_, err = db.Exec("UPDATE POST SET UpdatedAt = ?, Content = ? WHERE 	Id = ? ", u.UpdatedAt, u.Content, ID)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
func (d *Post) DeletePost(Id string) error {
	db, err := durable.CreateDatabase("./Database/", "common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM  POST WHERE  Id= ?", Id)

	d.UpdatedAt = time.Now().Unix()

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
