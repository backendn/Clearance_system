package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/backendn/clearance_system/db/sqlc"

	_ "github.com/lib/pq"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open("postgres",
		"postgres://root:secret@localhost:5433/university_clearance_test?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}

	testQueries = db.New(testDB)

	os.Exit(m.Run())
}
