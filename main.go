package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ozan1338/api"
	db "github.com/ozan1338/db/sqlc"
)

const (
	dbDriver = "mysql"
	dbSource = "root:root@tcp(localhost:3307)/simple_bank?parseTime=true"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	var err error
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("can't connect db: ", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ",err)
	}
}