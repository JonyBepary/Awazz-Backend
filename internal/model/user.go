package model

import (
	"log"
	"time"

	//"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (s *User) SaveUserData() error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS USER (
			UserName VARCHAR(128) PRIMARY KEY ,
			Password VARCHAR(128) NOT NULL,
			Email VARCHAR(128) ,
			ProfilePicUrl TEXT ,
			AccountTime INTEGER,
			BirthDate TEXT ,
			Gender TEXT ,
			LastEdit INTEGER)
	`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}

	statement, err := db.Prepare("INSERT INTO USER (UserName,Password,Email,ProfilePicUrl,AccountTime,BirthDate,Gender,LastEdit )VALUES(?,?,?,?,?,?,?,?)")

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec(s.UserName, s.Password, s.Email, s.ProfilePicUrl, s.AccountTime, s.BirthDate, s.Gender, s.LastEdit)
	if err != nil {
		panic(err)
	}
	return nil
}

func (g *User) GetUserData(userName string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from USER where UserName = ?", userName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&g.UserName, &g.Password, &g.Email, &g.ProfilePicUrl, &g.AccountTime, &g.BirthDate, &g.Gender, &g.LastEdit)

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
// previous data first before calling this
func (u *User) UpdatedUserData(UserName string) error {


	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	u.LastEdit = time.Now().Unix()

	_, err = db.Exec("UPDATE USER SET  LastEdit= ?, Password = ?, Email=?, BirthDate = ? WHERE 	UserName = ? ", u.LastEdit, u.Password, u.Email, u.BirthDate, UserName)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (d *User) DeleteUserData(UserName string) error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM  USER WHERE  UserName= ?", UserName)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
