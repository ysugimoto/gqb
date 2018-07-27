package gqb_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ysugimoto/gqb"
)

func connectDatabase() (*sql.DB, error) {
	return sql.Open("mysql", "root:root@tcp(127.0.0.1:"+os.Getenv("GQB_MYSQL_PORT")+")/gqb_test")
}

func BenchmarkNativeSQL(b *testing.B) {
	db, err := connectDatabase()
	if err != nil {
		b.Errorf("couldnt' connect database: %s", err.Error())
		return
	}
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := "SELECT * FROM companies WHERE id = ? OR id = ?"
		rows, err := db.Query(query, 1, 2)
		if err != nil {
			b.Errorf("failed to execute query on %d time: %s", i, err.Error())
			return
		}
		var cnt = 0
		for rows.Next() {
			var id int
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				b.Errorf("unexpected scan error: %s", err.Error())
				rows.Close()
				return
			}
			cnt++
		}
		rows.Close()
		if cnt != 2 {
			b.Errorf("unexpected result count. expect 2, actual %d", cnt)
			return
		}
	}
}

func BenchmarkQueryBuilder(b *testing.B) {
	db, err := connectDatabase()
	if err != nil {
		b.Errorf("couldnt' connect database: %s", err.Error())
		return
	}
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := gqb.New(db).
			Where("id", 1, gqb.Equal).
			OrWhere("id", 2, gqb.Equal).
			Get("companies")
		if err != nil {
			b.Errorf("failed to execute query on %d time: %s", i, err.Error())
			return
		}
		var cnt = 0
		for range rows {
			cnt++
		}
		if cnt != 2 {
			b.Errorf("unexpected result count. expect 2, actual %d", cnt)
			return
		}
	}
}
