package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

var (
	db         *sql.DB
	pgUsername string
	pgPassword string
	pgHost     string
	pgPort     string
	pgDBName   string
)

func init() {
	gotenv.Load()

	pgUsername = os.Getenv("POSTGRES_USERNAME")
	pgPassword = os.Getenv("POSTGRES_PASSWORD")
	pgHost = os.Getenv("POSTGRES_HOST")
	pgPort = os.Getenv("POSTGRES_PORT")
	pgDBName = os.Getenv("POSTGRES_DBNAME")
}

// connection fails fatally or returns a *sql.DB
func connection() *sql.DB {
	if db != nil {
		err := db.Ping()
		if err == nil {
			return db
		}
	}

	connstring := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		pgUsername,
		pgPassword,
		pgHost,
		pgPort,
		pgDBName,
	)

	dbconn, err = sql.Open("postgres", connstring)
	if err != nil {
		log.Fatalf("could not connect to db: %s\n", err.Error())
	}
	db = dbconn // cache for later calls
	return db
}
