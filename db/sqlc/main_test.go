package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open("mysql", "root:root@tcp(localhost:3307)/simple_bank?parseTime=true")

	if err != nil {
		log.Fatal("can't connect db: ", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
