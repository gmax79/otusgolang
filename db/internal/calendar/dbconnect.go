package calendar

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type dbConnect struct {
	db      *sql.DB
	notable bool
}

const createSchema = `
create table if not exists events (
id serial primary key,
alarm date not null,
information varchar(255) not null
);
`

const checkSchema = `
select column_name,data_type 
from information_schema.columns 
where table_name = 'events2' order by column_name;
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

func (c *dbConnect) checkSchema() error {
	db := c.db
	rows, err := db.Query(checkSchema)
	if err != nil {
		return err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	schemaErr := fmt.Errorf("Database schema different from required")
	if len(columns) != 2 {
		return schemaErr
	}
	c.notable = true
	tt := make(map[string]string)
	tt["alert"] = "date"
	tt["information"] = "character varying"
	tt["id"] = "integer"

	for rows.Next() {
		c.notable = false
		var cname, ctype string
		if err := rows.Scan(&cname, &ctype); err != nil {
			return err
		}
		if v, ok := tt[cname]; !ok {
			return schemaErr
		} else {
			if v != ctype {
				return schemaErr
			}
		}
	}
	return nil
}

func (c *dbConnect) CreateSchema() error {
	err := c.checkSchema()
	if err != nil {
		return err
	}
	if c.notable {
		_, err = c.db.Exec(createSchema)
	}
	return err
}

func (c *dbConnect) FindEvent() error {
	_, err := c.db.Exec("insert into events ")
	log(err)
	return err
}
