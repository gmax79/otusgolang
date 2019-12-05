package calendar

import (
	"database/sql"
	"fmt"
)

type dbSchema struct {
}

const createSchema = `
CREATE TABLE IF NOT EXISTS events (
timerid int NOT NULL,
information VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS timers (
id serial PRIMARY KEY,
alarm date NOT NULL
);
`

// DbSchemaError - type of error, where table schema not equal at checking
type DbSchemaError struct {
	tableName string
}

func (e *DbSchemaError) Error() string {
	return fmt.Sprintf("Table's %s schema is different from required", e.tableName)
}

// DbTableMissingError - table is missing in database
type DbTableMissingError struct {
	tableName string
}

func (e *DbTableMissingError) Error() string {
	return fmt.Sprintf("Table %s is missing in database", e.tableName)
}

const getTableSchema = `
select column_name,data_type 
from information_schema.columns 
where table_name = $1 order by column_name;
`

func (h dbSchema) checkTable(dbc *sql.DB, name string, schema map[string]string) error {
	var err error
	var rows *sql.Rows
	if rows, err = dbc.Query(getTableSchema, name); err != nil {
		return err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var cname, ctype string
		if err := rows.Scan(&cname, &ctype); err != nil {
			return err
		}
		v, ok := schema[cname]
		if !ok || v != ctype {
			return &DbSchemaError{tableName: name}
		}
		count++
	}
	if count == 0 {
		return &DbTableMissingError{tableName: name}
	}
	if count != len(schema) {
		return &DbSchemaError{tableName: name}
	}
	return nil
}

// CheckOrCreateSchema - function to create schema in empty db or error is schema is different
func (h dbSchema) CheckOrCreateSchema(dbc *sql.DB) error {
	et := map[string]string{
		"timerid":     "integer",
		"information": "character varying",
	}
	if err := skipMissedTable(h.checkTable(dbc, "events", et)); err != nil {
		return err
	}
	tt := map[string]string{
		"id":    "integer",
		"alarm": "date",
	}
	if err := skipMissedTable(h.checkTable(dbc, "timers", tt)); err != nil {
		return err
	}
	_, err := dbc.Exec(createSchema)
	return err
}

func skipMissedTable(err error) error {
	switch err.(type) {
	case *DbTableMissingError:
		return nil
	}
	return err
}
