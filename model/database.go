package model

import (
	"database/sql"
	"fmt"
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
	if db, err = sql.Open("sqlite3", dbpath); err != nil {
		log.Printf("error opening database: %v\n", err)
		defer Close()
		return err
	}

	// try to ping the database; this actually establishes the connection
	if err = db.Ping(); err != nil {
		log.Printf("error connecting to database: %v\n", err)
		defer Close()
		return err
	}

	// create the necessary tables
	statement := `create table if not exists product (id string not null primary key, description text);`
	if _, err = db.Exec(statement); err != nil {
		log.Printf("error executing %s: %v\n", statement, err)
		defer Close()
		return err
	}
	statement = `create table if not exists version (id string not null, product_id string not null, description text, primary key (product_id, id), foreign key(product_id) references product(id) on delete cascade on update no action);`
	if _, err = db.Exec(statement); err != nil {
		log.Printf("error executing %s: %v\n", statement, err)
		defer Close()
		return err
	}
	statement = `create table if not exists deployments (order int not null, product_id string not null, version_id string not null, description text, primary key(product_id, version_id, order), foreign key(product_id) references product(id) on delete cascade on update no action, foreign key (version_id) references verrsion(id) on delete cascade on update no action);`
	if _, err = db.Exec(statement); err != nil {
		log.Printf("error executing %s: %v\n", statement, err)
		defer Close()
		return err
	}
	return nil
}

// Load assumes that a database exists already under the given path.
func Load(dbpath string) error {
	var err error

	if _, err = os.Stat(dbpath); os.IsNotExist(err) {
		return err
	}

	if db, err = sql.Open("sqlite3", dbpath); err != nil {
		defer Close()
		log.Printf("error opening database: %v\n", err)
		return err
	}

	if err = db.Ping(); err != nil {
		defer Close()
		log.Printf("error connecing to database: %v\n", err)
		return err
	}

	return nil
}

// Close closes the database.
func Close() error {
	if db != nil {
		err := db.Close()
		db = nil
		return err
	}
	return nil
}

// GetAllProducts return the list of all registered products.
func GetAllProducts() ([]Product, error) {

	var err error

	if db == nil {
		log.Println("database is not open")
		return nil, fmt.Errorf("database is not open")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("error opening transaction: %v\n", err)
		return nil, err
	}

	rows, err := tx.Query("select id, description from products")
	if err != nil {
		log.Printf("error running query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	if err = tx.Commit(); err != nil {
		log.Printf("error committing transaction: %v\n", err)
		return nil, err
	}

	products := make([]Product, 0, 32)
	for rows.Next() {
		product := Product{}
		if err = rows.Scan(&product.id, &product.description); err != nil {
			log.Printf("error querying product: %v\n", err)
			continue
		}
		products = products(append, product)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error querying database: %v\n", err)
		return nil, err
	}

	return products, nil
}

// TODO: implement
/*
func GetProductByID(productId string) (Product, error) {
	var err error

	if db == nil {
		log.Println("database is not open")
		return nil, fmt.Errorf("database is not open")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("error opening transaction: %v\n", err)
		return nil, err
	}

	stmt, err = tx.Prepare("select id, descrption from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	rows, err := tx.Query("select id, description from products")
	if err != nil {
		log.Printf("error running query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	if err = tx.Commit(); err != nil {
		log.Printf("error committing transaction: %v\n", err)
		return nil, err
	}


}
*/
