package model

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

var (
	db         *sqlx.DB
	pgUser     string
	pgPassword string
	pgHostname string
	pgPort     string
	pgDB       string
	pgSSLMode  string
)

func init() {
	gotenv.Load()

	pgUser = os.Getenv("POSTGRES_USER")
	pgPassword = os.Getenv("POSTGRES_PASSWORD")
	pgHostname = os.Getenv("POSTGRES_HOSTNAME")
	pgPort = os.Getenv("POSTGRES_PORT")
	pgDB = os.Getenv("POSTGRES_DB")
	pgSSLMode = os.Getenv("POSTGRES_SSLMODE")
}

// database tries to return an open database
func database() (*sqlx.DB, error) {
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
		pgUser,
		url.QueryEscape(pgPassword),
		pgHostname,
		pgPort,
		pgDB,
		pgSSLMode,
	)

	log.Printf("connecting with connstring: %v\n",
		strings.Replace(connstring, pgPassword, "*", 1))

	dbconn, err := sqlx.Open("postgres", connstring)
	if err == nil {
		log.Printf("connected successfully, caching connection\n")
		db = dbconn // cache for later calls
	}
	log.Printf("going to return: (*sql.DB: %v, error: %v)", db, err)
	return db, err
}
