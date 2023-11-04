package durable

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	_ "github.com/glebarez/go-sqlite"
	"github.com/syndtr/goleveldb/leveldb"
)

// create a database connection to a file path
func CreateDatabase(params ...string) (*sql.DB, error) {

	// check if directory exists, create it if it doesn't
	dir := path.Join(params[0])
	dir = path.Join(dir, params[1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	if len(params) > 1 {
		dir = path.Join(dir, params[2])
	}
	fmt.Println(dir)
	db, err := sql.Open("sqlite", dir)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
func LevelDBCreateDatabase(params ...string) (*leveldb.DB, error) {

	// check if directory exists, create it if it doesn't
	dir := path.Join(params[0])
	dir = path.Join(dir, params[1])
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	if len(params) > 1 {
		dir = path.Join(dir, params[2])
	}
	fmt.Println(dir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}

	return db, err
}
