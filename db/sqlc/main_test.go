package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ozan1338/util"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	// var err error
	_,err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load configuration file: ",err)
	}

	testDb, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/simple_bank?parseTime=true")

	if err != nil {
		log.Fatal("can't connect db: ", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
