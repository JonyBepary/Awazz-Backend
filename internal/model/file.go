package model

import (
	"fmt"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (s *File) Save() error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `CREATE TABLE IF NOT EXISTS FILE(
		ULID VARCHAR(255) PRIMARY KEY,
		Name TEXT NOT NULL,
		CreatedAt INTEGER NOT NULL,
		UpdatedAt INTEGER,
		Hash TEXT NOT NULL,
		HashType TEXT NOT NULL,
		MimeType TEXT NOT NULL,
		Ext TEXT NOT NULL
	)`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}

	statement, err := db.Prepare("INSERT INTO FILE (ULID,Name,CreatedAt,UpdatedAt,Hash,HashType,MimeType,Ext) VALUES (?,?,?,?,?,?,?,?)")

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec(s.Uuid, s.Name, s.CreatedAt, s.UpdatedAt, s.Hash, s.HashType, s.MimeType, s.Ext)
	if err != nil {
		panic(err)
	}
	return nil
}

func (g *File) Get(ULID string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * from FILE where ULID = ?", ULID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&g.Uuid, &g.Name, &g.CreatedAt, &g.UpdatedAt, &g.Hash, &g.HashType, &g.MimeType, &g.Ext)
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

func (d *File) Delete(ULID string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM  FILE WHERE  ULID= ?", ULID)

	if err != nil {
		panic(err)
	}
	return nil
}

func (u *FileList) Get(ULIDS []string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, ULID := range ULIDS {
		var g File
		rows, err := db.Query(fmt.Sprintf("SELECT * from FILE where ULID = %v", ULID))
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		rows.Next()
		err = rows.Scan(g.Uuid, g.Name, g.CreatedAt, g.UpdatedAt, g.Hash, g.HashType, g.MimeType, g.Ext)
		if err != nil {
			panic(err)
		}
		u.Files = append(u.Files, &g)
		err = rows.Err()
		if err != nil {
			panic(err)
		}
	}

	return nil
}
