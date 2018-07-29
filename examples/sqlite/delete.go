package main

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ysugimoto/gqb"
)

func main() {
	db, err := sql.Open("sqlite3", "file:/tmp/gqb_test.sqlite?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	gqb.SetDriver("sqlite")
	_, err = gqb.New(db).
		Where("id", 3, gqb.Equal).
		Delete("companies")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Id: 3 deleted.")
}
