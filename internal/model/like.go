package model

import (
	"fmt"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (p *Like) SaveLikes(frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sql_cmd := `CREATE  TABLE IF NOT EXISTS LIKES (
		EntityId VARCHAR(255) PRIMARY KEY,
		EntityType VARCHAR(255),
		UserId VARCHAR(255),
		CreatedAt INTEGER,
		Unique (EntityId,UserId)
	)`

	_, err = db.Exec(sql_cmd)
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare("INSERT OR REPLACE INTO LIKES (EntityId,UserId,CreatedAt,EntityType) VALUES (?,?,?,?)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(p.EntityId, p.UserId, p.CreatedAt, &p.EntityType)
	if err != nil {
		panic(err)
	}

	return nil
}

func (l *Like) Delete(frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	delQuery := "DELETE FROM LIKES WHERE EntityId = ? AND UserId = ?"
	statement, err := db.Prepare(delQuery)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(l.EntityId, l.UserId)
	if err != nil {
		panic(err)
	}

	return nil
}

func (l *Likes) GetByEntityId(EntityId string, frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	getQuery := "SELECT * FROM LIKES WHERE EntityId = ?"
	rows, err := db.Query(getQuery, EntityId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var like Like
		err = rows.Scan(&like.EntityId, &like.UserId, &like.CreatedAt, &like.EntityType)
		if err != nil {
			panic(err)
		}
		l.Likes = append(l.Likes, &like)
	}
	return nil
}

func (l *Likes) GetByUserId(UserId string, frag_num int64) error {
	db, err := durable.CreateDatabase("Database/", "Common", fmt.Sprintf("Shard_%d.sqlite", frag_num))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	getQuery := "SELECT * FROM LIKES WHERE UserId = ?"
	rows, err := db.Query(getQuery, UserId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var like Like
		err = rows.Scan(&like.EntityId, &like.UserId, &like.CreatedAt, &like.EntityType)
		if err != nil {
			panic(err)
		}
		l.Likes = append(l.Likes, &like)
	}

	return nil
}
