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

func (p *dbProvder) AddTrigger(d date) (string, error) {
	timer := d.String()
	request := `INSERT INTO timers (timer) 
	SELECT $1::timestamp WHERE NOT EXISTS
	(SELECT 1 FROM timers WHERE timer = $1::timestamp)`
	_, err := p.db.Exec(request, timer)
	if err != nil {
		return "", err
	}
	return timer, nil
}

func (p *dbProvder) AddEvent(d date, info string) error {
	timer := d.String()
	_, err := p.db.Exec("INSERT INTO events timer, information values($1, $2)", timer, info)
	return err
}

func (p *dbProvder) GetTriggers() ([]Date, error) {
	rows, err := p.db.Query("SELECT timer FROM timers")
	if err != nil {
		return []Date{}, err
	}
	defer rows.Close()
	ids := make([]Date, 0, 10)
	for rows.Next() {
		var timer string
		err := rows.Scan(&timer)
		if err != nil {
			return []Date{}, err
		}
		var d date
		d.ParseDate(timer)
		ids = append(ids, d.d)
	}
	return ids, nil
}

func (p *dbProvder) DeleteTrigger(d date) error {
	timer := d.String()
	_, err := p.db.Exec("DELETE FROM events WHERE timer = $1; ", timer)
	if err != nil {
		return err
	}
	_, err = p.db.Exec("DELETE FROM triggers WHERE timer = $1; ", timer)
	return err
}

func (p *dbProvder) Invoke(id string) {
	//	p.db.Query("SELECT ")
}
