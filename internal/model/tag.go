package model

import (
	"log"
	"runtime"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
)

func (t *Tags) SetTag() error {
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("Error in %s:%d: %s", file, line, err.Error())
		return err
	}
	defer db.Close()
	sql_cmd := `CREATE TABLE IF NOT EXISTS TAGS (Id VARCHAR(255),
    Type VARCHAR(255),
    Tag VARCHAR(255),
    UNIQUE(Id, Type, Tag)
	)`

	_, err = db.Exec(sql_cmd)
	if err != nil{
		_, file, line, _ := runtime.Caller(0)
		log.Printf("Error in %s:%d: %s", file, line, err.Error())
		return err
	}

	statement, err := db.Prepare("INSERT INTO TAGS (Id,Type,Tag) VALUES (?,?,?)")
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("Error in %s:%d: %s", file, line, err.Error())
		return err
	}
	_, err = statement.Exec(t.Id, t.Type, t.Tag)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Printf("Error in %s:%d: %s", file, line, err.Error())
		return err
	}

	return nil
}
