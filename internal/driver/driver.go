package driver

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDBLifetime = 5 * time.Minute

// ConnectSQL creates database pool for PostgreSQL
func ConnectSQL(dsn string) (*DB, error) {

	d, err := newDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = d

	return dbConn, nil
}

// newDatabase creates a new database for the application
func newDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
