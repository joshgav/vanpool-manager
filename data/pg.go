package data

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

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
	pgSSLMode  string
)

func init() {
	gotenv.Load()

	pgUsername = os.Getenv("POSTGRES_USER")
	pgPassword = os.Getenv("POSTGRES_PASSWORD")
	pgHost = os.Getenv("POSTGRES_HOSTNAME")
	pgPort = os.Getenv("POSTGRES_PORT")
	pgDBName = os.Getenv("POSTGRES_DB")
	pgSSLMode = os.Getenv("POSTGRES_SSLMODE")
}

// connection tries to return an open database
func connection() (*sql.DB, error) {
	if db != nil {
		log.Printf("found cached db %v\n", db)
		err := db.Ping()
		if err == nil {
			log.Printf("ping succeeded\n")
			return db, err
		}
		log.Printf("cached db failed Ping, error: %v\n", err)
		log.Printf("will try from scratch\n")
	}

	connstring := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pgUsername,
		url.QueryEscape(pgPassword),
		pgHost,
		pgPort,
		pgDBName,
		pgSSLMode,
	)

	log.Printf("connecting with connstring: %v\n",
		strings.Replace(connstring, pgPassword, "*", 1))

	dbconn, err := sql.Open("postgres", connstring)
	if err == nil {
		log.Printf("connected successfully, caching connection\n")
		db = dbconn // cache for later calls
	}
	log.Printf("going to return: (*sql.DB: %v, error: %v)", db, err)
	return db, err
}
