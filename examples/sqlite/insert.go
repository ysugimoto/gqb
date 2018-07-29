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
	data := gqb.Data{"name": "Slack"}
	// the result is sql.Result
	result, err := gqb.New(db).Insert("companies", data)

	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to retrieve last inserted ID")
		return
	}
	fmt.Printf("Id %d has been inserted", id)
}
