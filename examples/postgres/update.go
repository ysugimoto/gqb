package main

import (
	"fmt"
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
	data := gqb.Data{"name": "Github"}
	_, err = gqb.New(db).
		Where("id", 1, gqb.Equal).
		Update("companies", data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Id: 1 updated to Github.")

}
