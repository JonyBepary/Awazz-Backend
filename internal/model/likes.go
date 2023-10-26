package model

import "github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"

func (p *Likes) SaveLikes() error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sql_cmd := `CREATE TABLE likes (
		PostId VARCHAR(255) PRIMARY KEY,
		UserId VARCHAR(255),
		CommentId VARCHAR(255),
		CreatedAt INTEGER,
	)`

	_, err = db.Exec(sql_cmd)
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare("INSERT INTO likes (PostId,UserId,CommentId,CreatedAt) VALUES (?,?,?,?)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.PostId, p.UserId, p.CommentId, p.CreatedAt)
	if err != nil {
		panic(err)
	}

	return nil

}

func (p *Likes) GetLikes(pid string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// spew.Dump(rows)
	return nil
}
