package main

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ysugimoto/gqb"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/example")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	gqb.SetDriver("mysql")
	data := gqb.Data{"name": "Slack"}
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
