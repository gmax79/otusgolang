package calendar

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type dbConnect struct {
	db *sql.DB
}

const createSchema = `
	create table if not exists events (
		id serial primary key,
		information  varchar(255) not null
	);
	create table if not exists triggers {
		id serial primary key,
		alert date
		
	}
`

// Connect - create connection
func dbconnect(pghost string) (*dbConnect, error) {
	connection, err := sql.Open("pgx", pghost) // *sql.DB
	if err != nil {
		return nil, err
	}
	return &dbConnect{db: connection}, nil
}

func log(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (c *dbConnect) CreateSchema() error {
	db := c.db
	_, err := db.Exec(createSchema)
	log(err)
	return err
}

func (c *dbConnect) FindEvent() error {
	_, err := c.db.Exec("insert into events ")
	log(err)
	return err
}
