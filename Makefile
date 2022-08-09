postgres:
	docker run --name mysql_docker -p 3307:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql

createdb:
	docker exec -it mysql_docker createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it mysql_docker dropdb --username=root --owner=root simple_bank

migrateup:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3307)/simple_bank" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:secret@tcp(localhost:3307)/simple_bank" -verbose down

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