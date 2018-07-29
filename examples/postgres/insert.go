package main

import (
	"log"

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

	gqb.SetDriver("postgres")
	data := gqb.Data{"name": "Slack"}
	// the result is sql.Result
	if _, err := gqb.New(db).Insert("companies", data); err != nil {
		log.Fatal(err)
	}
}
