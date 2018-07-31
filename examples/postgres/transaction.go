package main

import (
	"log"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/ysugimoto/gqb"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	gqb.SetDriver("postgres")
	data := gqb.Data{
		"name":       "Slack",
		"created_at": gqb.Datetime(time.Now()),
	}
	// gqb.New() also accepts *sql.Tx
	if _, err := gqb.New(tx).Insert("companies", data); err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	tx.Commit()
}
