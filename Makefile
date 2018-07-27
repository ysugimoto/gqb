.PHONY: test bench

test:
	go test .

ci: test
	./scripts/generate-my-conf.sh
	./scripts/wait-for-database.sh
	./scripts/create-example-data.sh
	go run examples/select.go
	go run examples/select_context.go
	go run examples/select_join.go
	go run examples/select_one.go
	go run examples/transaction.go
	go run examples/insert.go
	go run examples/update.go
	go run examples/delete.go

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
