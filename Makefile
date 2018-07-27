.PHONY: test bench

test:
	go test .

bench:
	./scripts/generate-my-conf.sh
	./scripts/wait-for-database.sh
	./scripts/create-test-data.sh 100
	go test -bench . -benchmem
	./scripts/create-test-data.sh 1000
	go test -bench . -benchmem
	./scripts/create-test-data.sh 10000
	go test -bench . -benchmem

local-bench:
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
