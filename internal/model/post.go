package model

import (
	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (p *Post) SavePost() error {
	db, err := durable.CreateDatabase("Database/post.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS Post (
		Id VARCHAR(128) PRIMARY KEY,
		Community VARCHAR(128),
		Content TEXT,
		CreatedAt TIMESTAMP,
		UpdatedAt TIMESTAMP,
		DeletedAt TIMESTAMP,
		Likes INT,
	    Shares INT,
	    Comments INT,
	    Author VARCHAR(128),
	    Parent INT,
	    Rank INT,
	    IsSensitive bool,
	    IsNsfw bool,
	    IsDeleted bool,
	    IsPinned bool,
	    IsEdited bool,
	    IsLiked bool,
	    IsShared bool,
	    IsCommented bool,
	    IsSubscribed bool,
	    IsBookmarked bool,
	    IsReblogged bool,
	    IsMentioned bool,
	    IsPoll bool,
	    IsPollVoted bool,
	    IsPollExpired bool,
	    IsPollClosed bool,
	    IsPollMultiple bool,
	    IsPollHideTotals bool)
	`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}
	str2 := `
INSERT INTO Post (Id,Community,Content,CreatedAt,UpdatedAt,DeletedAt,Likes,Shares,Comments,Author,Parent,Rank,IsSensitive,IsNsfw,IsDeleted,IsPinned,IsEdited,IsLiked,IsShared,IsCommented,IsSubscribed,IsBookmarked,IsReblogged,IsMentioned,IsPoll,IsPollVoted,IsPollExpired,IsPollClosed,IsPollMultiple,IsPollHideTotals) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`
	statement, err := db.Prepare(str2)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.Id,
		p.Community,
		p.Content,
		p.CreatedAt.String(),
		p.UpdatedAt.String(),
		p.DeletedAt.String(),
		p.Likes,
		p.Shares,
		p.Comments,
		p.Author,
		p.Parent,
		p.Rank,
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
		p.IsPollHideTotals,
	)
	if err != nil {
		panic(err)
	}

	return nil
}

/*func (p *Post) GetPost() error {
	db,
		err := durable.CreateDatabase("Database/post")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return nil
}
func (p *Post) UpdatePost() error {
	db, err := durable.CreateDatabase("Database/Post")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return nil
}
func (p *Post) DeletePost() error {
	db, err := durable.CreateDatabase("Database/Post")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return nil
}
*/
