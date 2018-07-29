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

	gqb.SetDriver("mysql")
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
