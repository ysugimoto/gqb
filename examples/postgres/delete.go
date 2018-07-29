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
	_, err = gqb.New(db).
		Where("id", 3, gqb.Equal).
		Delete("companies")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Id: 3 deleted.")
}
