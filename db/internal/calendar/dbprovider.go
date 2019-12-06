package calendar

import (
	"database/sql"

	_ "github.com/jackc/pgx/stdlib" // attach pgx postgres driver
)

func connectToDatabase(dsn string) (*sql.DB, error) {
	connection, err := sql.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, err
	}
	return connection, nil
}

type dbProvder struct {
	db *sql.DB
}

func getProvider(db *sql.DB) *dbProvder {
	return &dbProvder{db: db}
}

func (p *dbProvder) AddTrigger(trigger string) (int64, error) {
	request := `INSERT INTO timers (alarm) 
	SELECT $1 WHERE NOT EXISTS (SELECT id FROM timers WHERE alarm = '$2')
	RETURNING id;`

	result, err := p.db.Exec(request, trigger, trigger)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
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
