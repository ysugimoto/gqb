package main

import (
	"fmt"
	"log"
	"time"

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

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	gqb.SetDriver("sqlite")
	data := gqb.Data{
		"name":       "Slack",
		"created_at": gqb.Datetime(time.Now()),
	}
	// gqb.New() also accepts *sql.Tx
	result, err := gqb.New(tx).Insert("companies", data)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Println("failed to retrieve last inserted ID")
		return
	}
	tx.Commit()
	fmt.Printf("Id %d has been inserted", id)
}
