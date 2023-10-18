package durable

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

// create a database connection to a file path
func CreateDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

func SetColumn(col, col_types []string) string {
	STR := ""
	for i := 0; i < len(col); i++ {
		STR += fmt.Sprintf("%v %v, ", col[i], col_types[i])
	}
	return STR[:len(STR)-2]
}
