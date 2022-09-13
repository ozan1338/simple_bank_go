mysql:
	docker run --name mysql_docker -p 3307:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql

mysqlnew:
	docker run --name simple_bank --network bank_network -p 8080:8080 -e GIN_MODE=release  -e DB_SOURCE="root:root@tcp(mysql_docker:3306)/simple_bank?parseTime=true" simplebank:latest

createdb:
	docker exec -it mysql_docker createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it mysql_docker dropdb --username=root --owner=root simple_bank

inspectContainer:
	docker container inspect nama_container

inspectNetwork:
	docker network inspect nama_network

createNetwork:
	docker network create nama_network

connectNetwork:
	docker network connect nama_container nama_network

migrateup:
	migrate -path db/migration -database "mysql://root:root@tcp(localhost:3307)/simple_bank" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:root@tcp(localhost:3307)/simple_bank" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./..

runbashmysql:
	docker exec -it mysql_docker /bin/bash

runmysql:
	mysql -u root -p

runmockgenreflectmode:
	mockgen --build_flags=--mod=mod -package mockdb  -destination db/mock/store.go github.com/ozan1338/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test