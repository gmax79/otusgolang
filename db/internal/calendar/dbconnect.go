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
		information  varchar(255) not null,
		alert date
	);
`

// Connect - create connection
func dbconnect(pghost string) (*dbConnect, error) {
	connection, err := sql.Open("pgx", pghost) // *sql.DB
	if err != nil {
		return nil, err
	}
	return &dbConnect{db: connection}, nil
}

func (c *dbConnect) CreateSchema() error {
	db := c.db
	_, err := db.Exec(createSchema)
	fmt.Println(err)
	return err
}

func (c *dbConnect) AddEvent() error {
	_, err := c.db.Exec("insert into events ")
	fmt.Println(err)
	return err
}
