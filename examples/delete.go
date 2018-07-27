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

	// Well, we have to remove it...
	_, err = gqb.New(db).
		Where("id", 4, gqb.Equal).
		Delete("companies")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Id: 4 deleted.")
}
