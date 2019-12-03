package calendar

type dbHandler struct {
	db *dbConnection
}

const createSchema = `
create table if not exists events (
triggerid int not null,
information varchar(255) not null
);
create table if not exists triggers (
id serial primary key,
alarm date not null,
);
`

// Connect - create connection
func dbconnect(pghost string) (*dbHandler, error) {
	connection, err := connectToDatabase(pghost)
	if err != nil {
		return nil, err
	}
	return &dbHandler{db: connection}, nil
}

func (c *dbHandler) CheckOrCreateSchema() error {
	et := map[string]string{
		"triggerid":   "integer",
		"information": "character varying",
	}
	if err := skipMissedTable(c.db.checkTable("events", et)); err != nil {
		return err
	}
	tt := map[string]string{
		"id":    "integer",
		"alarm": "date",
	}
	if err := skipMissedTable(c.db.checkTable("triggers", tt)); err != nil {
		return err
	}
	err := c.db.Exec(createSchema)
	return err
}

/*func (c *dbHandler) AddEvent(d Date, info string) error {
	c.db.Exec("insert into events date, information values($1, $2)", d.String(), info)
}*/

/*func (c *dbHandler) FindEvent(d Date) error {
	rows, err := c.db.Query("select id from events where alarm=$1", d.String())
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
	}
	return err
}*/
