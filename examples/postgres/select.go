package main

import (
	"fmt"
	"log"

	"database/sql"
	"encoding/json"

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
	results, err := gqb.New(db).
		Select("name").
		Where("id", 1, gqb.Equal).
		Get("companies")

	if err != nil {
		log.Fatal(err)
	}
	// retrieve result
	for _, r := range results {
		fmt.Println(r.MustString("name")) //=> Google
	}

	// Also can marshal JSON directly
	buf, _ := json.Marshal(results)
	fmt.Println(string(buf)) //=> [{"name":"Google"}]

	// Map to your struct
	type Company struct {
		Name string `db:"name"` // gqb maps value corresponds to "db" tag field
	}
	companies := []Company{}
	if err := results.Map(&companies); err != nil {
		log.Fatal(err)
	}
	fmt.Println(companies[0].Name) //=> Google
}
