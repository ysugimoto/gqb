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

	gqb.SetDriver("mysql")
	result, err := gqb.New(db).
		Select("name").
		Where("id", 1, gqb.Equal).
		GetOne("companies")

	if err != nil {
		log.Fatal(err)
	}
	// retrieve result
	fmt.Println(result.MustString("name")) //=> Google

	// Also can marshal JSON directly
	buf, _ := json.Marshal(result)
	fmt.Println(string(buf)) //=> {"name":"Google"}

	// Map to your struct
	type Company struct {
		Name string `db:"name"` // gqb maps value corresponds to "db" tag field
	}
	company := Company{}
	if err := result.Map(&company); err != nil {
		log.Fatal(err)
	}
	fmt.Println(company.Name) //=> Google
}
