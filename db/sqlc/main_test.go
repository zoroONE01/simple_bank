package db

import (
	"database/sql"
	"log"
	"os"
	"simple_bank/utils"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB
var testStore Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	testDB = conn
	testQueries = New(testDB)
	testStore = NewStore(testDB)
	os.Exit(m.Run())
}
