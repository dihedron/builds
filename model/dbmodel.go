package model

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// The private database copy.
var db *sql.DB

// New creates a new database, deleting any pre-existing instance.
func New(dbpath string) error {

	var err error

	// close the db if already open
	if db != nil {
		db.Close()
	}
	// check if a file already exists, and if so remove it
	if stat, err := os.Stat(dbpath); err != nil && !stat.IsDir() {
		os.Remove(dbpath)
	}
	// open the new database
	db, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		return err
	}
	// create the necessary tables
	statement := `create table if not exists product (id string not null primary key, description text);`
	if _, err = db.Exec(statement); err != nil {
		log.Printf("error executing %s: %v\n", statement, err)
		return err
	}
	statement = `create table if not exists version (id string not null primary key, description text, foreign key(product_id) references product(id));`
	if _, err = db.Exec(statement); err != nil {
		log.Printf("error executing %s: %v\n", statement, err)
		return err
	}
	statement = `create table if not exists deployments (id string not null primary key, description text, foreign key(product_id) references product(id));`
	if _, err = db.Exec(statement); err != nil {
		log.Printf("error executing %s: %v\n", statement, err)
		return err
	}
	return nil
}

// Load assumes that a database exists already under the given path.
func Load(dbpath string) error {
	var err error

	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		return err
	}

	db, err = sql.Open("sqlite3", dbpath)
	if err != nil {
		return err
	}

	return nil
}

func SaveToFile(dbpath string) error {
	return nil
}
