package model

import (
	"fmt"
	"log"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (p *Post) SavePost(frag_num int64) error {
	db, err := durable.CreateDatabase("./Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	str := `
	CREATE TABLE IF NOT EXISTS POST(
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

	str2 := `INSERT INTO POST (Id,Community,Content,CreatedAt,UpdatedAt,DeletedAt,Likes,Shares,Comments,Author,Parent,Rank,IsSensitive,IsNsfw,IsDeleted,IsPinned,IsEdited,IsLiked,IsShared,IsCommented,IsSubscribed,IsBookmarked,IsReblogged,IsMentioned,IsPoll,IsPollVoted,IsPollExpired,IsPollClosed,IsPollMultiple,IsPollHideTotals,FragmentationKey) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	statement, err := db.Prepare(str2)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Id, p.Community, p.Content, p.CreatedAt, p.UpdatedAt, p.DeletedAt, p.Likes, p.Shares, p.Comments, p.Author, p.Parent, p.Rank, p.IsSensitive, p.IsNsfw, p.IsDeleted, p.IsPinned, p.IsEdited, p.IsLiked, p.IsShared, p.IsCommented, p.IsSubscribed, p.IsBookmarked, p.IsReblogged, p.IsMentioned, p.IsPoll, p.IsPollVoted, p.IsPollExpired, p.IsPollClosed, p.IsPollMultiple, p.IsPollHideTotals, p.FragmentationKey)
	if err != nil {
		panic(err)
	}

	return nil
}
func (p *Post) GetPost(msgId string, frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// spew.Dump(rows)
	rows, err := db.Query("SELECT * FROM POST WHERE Id = ?", msgId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(p.Id, p.Community, p.Content, p.CreatedAt, p.UpdatedAt, p.DeletedAt, p.Likes, p.Shares, p.Comments, p.Author, p.Parent, p.Rank, p.IsSensitive, p.IsNsfw, p.IsDeleted, p.IsPinned, p.IsEdited, p.IsLiked, p.IsShared, p.IsCommented, p.IsSubscribed, p.IsBookmarked, p.IsReblogged, p.IsMentioned, p.IsPoll, p.IsPollVoted, p.IsPollExpired, p.IsPollClosed, p.IsPollMultiple, p.IsPollHideTotals, p.FragmentationKey)
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
func (u *Post) UpdatePost(ID string, frag_num int64) error {
	db, err := durable.CreateDatabase("./Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
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
func (d *Post) DeletePost(Id string, frag_num int64) error {
	db, err := durable.CreateDatabase("./Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM  POST WHERE  Id= ?", Id)
	if err != nil {
		panic(err)
	}

	d.UpdatedAt = time.Now().Unix()

	return nil
}
