package main

import (
	"fmt"
	"log"
	"time"

	"database/sql"
	"encoding/json"

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
	results, err := gqb.New(db).
		Select("name", "created_at").
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
		Name      string    `db:"name"` // gqb maps value corresponds to "db" tag field
		CreatedAt time.Time `db:"created_at"`
	}
	companies := []Company{}
	if err := results.Map(&companies); err != nil {
		log.Fatal(err)
	}
	fmt.Println(companies[0].Name) //=> Google
	fmt.Println(companies[0].CreatedAt.Format("2006-01-02 15:04:05"))
}
