package model

import (
	"fmt"
	"log"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

/*
	CRUD FOR FOLLOWER
*/

func (s *Follower) SaveFollower() error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS Follower (
		UserId VARCHAR(255) PRIMARY KEY,
		Status  BOOL,
		Time  INTEGER,
		FollowAccount VARCHAR(255),
		UnfollowAccount VARCHAR(255)
	)

		`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}

	statement, err := db.Prepare("INSERT INTO Follower (UserId,Status,Time,FollowAccount,UnfollowAccount) VALUES (?,?,?,?,?)")

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec(s.UserId, s.Status, s.Time, s.FollowAccount, s.UnfollowAccount)
	if err != nil {
		panic(err)
	}
	return nil
}

func (g *Follower) GetFollower(UserID string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select * from Follower where UserID = %v", UserID))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&g.UserId, &g.Status, &g.Time, &g.FollowAccount, &g.UnfollowAccount)

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

// func (u *Follower) UpdatedFollower(UserId string) error {
// 	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// 	u.Time = time.Now().Unix()
// 	_, err = db.Exec("UPDATE Follower SET  Time= ?,status = ?, FollowAccount=?, unfollowAccount=?  WHERE msId= ?", u.Time, u.Status, u.FollowAccount, u.UnfollowAccount, UserId)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return nil
// }

func (d *Follower) DeleteFollowee(FolloweeId string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// another rule don't type something that can be copied
	// typing on your own tends to bring more mistake
	// baki table check koro okay): abar thik ache bhai
	_, err = db.Exec("DELETE FROM  Follower WHERE  UserId= ?", FolloweeId) /// mair khabatable name thik koroe

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

/*
	CRUD FOR Followee
*/

func (s *Followee) SaveFollowee() error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	str := `
	CREATE TABLE IF NOT EXISTS Followee (
		UserId VARCHAR(255) PRIMARY KEY,
		BlockAccountLink TEXT,
		BlockAccountName TEXT
	)
		`
	_, err = db.Exec(str)
	if err != nil {
		panic(err)
	}

	statement, err := db.Prepare("INSERT INTO Followee (UserId,BlockAccountLink,BlockAccountName) VALUES (?,?,?)")

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec(s.UserId, s.BlockAccountLink, s.BlockAccountLink, s.BlockAccountName)
	if err != nil {
		panic(err)
	}
	return nil
}

func (g *Followee) GetFollowee(UserId string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("select * from Followee where UserID = %v", UserId))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&g.UserId, &g.BlockAccountLink, &g.BlockAccountName)

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

// func (u *Followee) UpdatedFollowee(UserId string) error {
// 	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
// 	u.Time = time.Now().Unix()
// 	_, err = db.Exec("UPDATE Followee SET  Time = ?,BlockAccountLink = ?, BlockAccountName=?WHERE UserId= ?", u.Time, u.BlockAccountLink, u.BlockAccountName, UserId)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return nil
// }

func (d *Followee) DeleteFollower(FollowerId string) error {
	db, err := durable.CreateDatabase("./Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM  Followee WHERE  UserId= ?", FollowerId)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
