package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ozan1338/api"
	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/util"
)

func main() {
	config,err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config: ",err)
	}

	conn, err := sql.Open(config.DBDrvier, config.DBSource)

	if err != nil {
		log.Fatal("can't connect db: ", err)
	}

	
	store := db.NewStore(conn)

	server,err := api.NewServer(config,store)

	if err != nil {
		log.Fatal("can't start server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ",err)
	}
}