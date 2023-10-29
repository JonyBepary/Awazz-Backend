package model

import (
	"log"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (s *Token) SaveToken() error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS TOKEN (
		UserName VARCHAR(128) PRIMARY KEY,
		Token  VARCHAR(128),
		GenerateTime INTEGER)
	`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}

	// check if user already exist

	statement, err := db.Prepare("INSERT OR REPLACE INTO TOKEN (UserName, Token, GenerateTime) VALUES (?,?,?)")
	if err != nil {
		panic(err)
	}

	_, err = statement.Exec(s.UserName, s.Token, s.GenerateTime)
	if err != nil {
		panic(err)
	}
	return nil
}

func (g *Token) GetTokenFromDB(UserName string) error {

	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from TOKEN where UserName = ?", UserName)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&g.UserName, &g.Token, &g.GenerateTime)

		if err != nil {
			log.Fatal(err)
		}

	}
	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (u *Token) UpdatedToken(UserName string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	u.GenerateTime = time.Now().Unix()
	_, err = db.Exec("UPDATE TOKEN SET  GenerateTime= ?, Token = ? WHERE UserName = ? ", u.GenerateTime, u.Token, UserName)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (d *Token) DeleteToken(UserName string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM  TOKEN WHERE  UserName= ?", UserName)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
