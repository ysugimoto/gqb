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

	_, err = gqb.New(db).
		Where("id", 3, gqb.Equal).
		Delete("companies")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Id: 3 deleted.")
}
