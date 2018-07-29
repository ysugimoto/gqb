.PHONY: test bench ci mysql postgres sqlite

test:
	go test .

e2e: mysql postgres sqlite

mysql:
	# MySQL test
	./scripts/mysql/generate-my-conf.sh
	./scripts/mysql/wait-for-database.sh
	./scripts/mysql/create-example-data.sh
	go run examples/mysql/select.go
	go run examples/mysql/select_context.go
	go run examples/mysql/select_join.go
	go run examples/mysql/select_one.go
	go run examples/mysql/transaction.go
	go run examples/mysql/insert.go
	go run examples/mysql/update.go
	go run examples/mysql/delete.go

postgres:
	# PostgreSQL test
	./scripts/postgres/wait-for-database.sh
	./scripts/postgres/create-example-data.sh
	go run examples/postgres/select.go
	go run examples/postgres/select_context.go
	go run examples/postgres/select_join.go
	go run examples/postgres/select_one.go
	go run examples/postgres/transaction.go
	go run examples/postgres/insert.go
	go run examples/postgres/update.go
	go run examples/postgres/delete.go

sqlite:
	# SQLite test
	./scripts/sqlite/create-example-data.sh
	go run examples/sqlite/select.go
	go run examples/sqlite/select_context.go
	go run examples/sqlite/select_join.go
	go run examples/sqlite/select_one.go
	go run examples/sqlite/transaction.go
	go run examples/sqlite/insert.go
	go run examples/sqlite/update.go
	go run examples/sqlite/delete.go

bench:
	docker ps | grep "gqb_mysql_test" | awk '{print $$1}' | xargs docker stop
	docker run --rm -d --name gqb_mysql_test -e "MYSQL_ROOT_PASSWORD=root" -p $(GQB_MYSQL_PORT):3306 mysql:5.7
	./scripts/generate-my-conf.sh
	./scripts/wait-for-database.sh
	./scripts/create-test-data.sh 100
	go test -bench . -benchmem
	./scripts/create-test-data.sh 1000
	go test -bench . -benchmem
	./scripts/create-test-data.sh 10000
	go test -bench . -benchmem
	docker stop gqb_mysql_test
