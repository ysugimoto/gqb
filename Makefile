.PHONY: test bench

test:
	go test .

bench:
	docker run --rm -d --name gqb_mysql_test -e "MYSQL_ROOT_PASSWORD=root" -p 63306:3306 mysql:5.7
	./scripts/wait-for-database.sh
	./scripts/create-test-data.sh
	go test -bench . -benchmem
	docker stop gqb_mysql_test
