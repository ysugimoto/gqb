package main

import (
	"fmt"
	"log"

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

	result, err := gqb.New(db).
		Select("name", "url").
		Join("company_attributes", "id", "company_id", gqb.Equal).
		// Need to specify table name to avoid ambigous 'id'column
		Where("companies.id", 3, gqb.Equal).
		GetOne("companies")

	if err != nil {
		log.Fatal(err)
	}
	// retrieve result
	fmt.Println(result.MustString("url")) //=> https://microsoft.com

	// Also can marshal JSON directly
	buf, _ := json.Marshal(result)
	fmt.Println(string(buf)) //=> {"name":"Microsoft","url":"https://microsoft.com"}

	// Map to your struct
	type Company struct {
		Name string `db:"name"`
		Url  string `db:"url"`
	}
	ms := Company{}
	if err := result.Map(&ms); err != nil {
		log.Fatal(err)
	}
	fmt.Println(ms.Name) //=> Microsoft
}
