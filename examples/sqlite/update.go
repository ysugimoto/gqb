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
	data := gqb.Data{"name": "Github"}
	_, err = gqb.New(db).
		Where("id", 1, gqb.Equal).
		Update("companies", data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Id: 1 updated to Github.")

}
