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
	config,err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load configuration file: ",err)
	}

	testDb, err = sql.Open(config.DBDrvier, config.DBSoureTesting)

	// if(config.DBLOCAL) {
	// 	testDb, err = sql.Open(config.DBDrvier, config.DBSource)
	// } else {
	// 	testDb, err = sql.Open(config.DBDrvier, config.DBSoureTesting)
	// }


	if err != nil {
		log.Fatal("can't connect db: ", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
