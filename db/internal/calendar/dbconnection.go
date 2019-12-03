package calendar

import (
	"database/sql"
	"fmt"

	// attach pgx postgres driver
	_ "github.com/jackc/pgx/stdlib"
)

type dbConnection struct {
	db *sql.DB
}

func connectToDatabase(dsn string) (*dbConnection, error) {
	connection, err := sql.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, err
	}
	return &dbConnection{db: connection}, nil
}

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

func (h *dbConnection) checkTable(name string, schema map[string]string) error {
	var err error
	var rows *sql.Rows
	if rows, err = h.db.Query(getTableSchema, name); err != nil {
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

func (h *dbConnection) Exec(request string, args ...interface{}) error {
	_, err := h.db.Exec(request, args...)
	return err
}

func (h *dbConnection) Query(request string, args ...interface{}) (*sql.Rows, error) {
	r, err := h.db.Query(request, args...)
	return r, err
}

func skipMissedTable(err error) error {
	switch err.(type) {
	case *DbTableMissingError:
		return nil
	}
	return err
}
