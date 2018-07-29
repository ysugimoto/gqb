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
	data := gqb.Data{"name": "Github"}
	_, err = gqb.New(db).
		Where("id", 1, gqb.Equal).
		Update("companies", data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Id: 1 updated to Github.")

}
