package calendar

import (
	"database/sql"
	"fmt"

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

func (p *dbProvder) GetTriggers() ([]Date, error) {
	rows, err := p.db.Query("SELECT DISTINCT timer FROM events;")
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
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (p *dbProvder) DeleteTrigger(d date) error {
	timer := d.String()
	_, err := p.db.Exec("DELETE FROM events WHERE timer = $1;", timer)
	return err
}

func (p *dbProvder) AddEvent(d date, info string) error {
	timer := d.String()
	_, err := p.db.Exec("INSERT INTO events (timer, information) VALUES($1, $2) ON CONFLICT (timer, information) DO NOTHING;", timer, info)
	return err
}

func (p *dbProvder) GetEventsCount(d date) (int, error) {
	timer := d.String()
	var count int
	err := p.db.QueryRow("SELECT COUNT (*) FROM events WHERE timer = $1", timer).Scan(&count)
	return count, err
}

func (p *dbProvder) DeleteEventIndex(d date, index int) error {
	request := "DELETE FROM events WHERE ctid IN (SELECT ctid FROM events WHERE timer = $1::timestamp limit 1 offset $2);"
	timer := d.String()
	_, err := p.db.Exec(request, timer, index)
	return err
}

func (p *dbProvder) DeleteEvent(d date, e Event) error {
	timer := d.String()
	_, err := p.db.Exec("DELETE FROM events WHERE timer = $1::timestamp AND information = $2;", timer, string(e))
	return err
}

func (p *dbProvder) GetEvent(d date, index int) (Event, error) {
	request := "SELECT information FROM events WHERE timer = $1::timestamp limit 1 offset $2;"
	timer := d.String()
	var e Event
	err := p.db.QueryRow(request, timer, index).Scan(&e)
	return e, err
}

func (p *dbProvder) MoveEvent(d date, e Event, to date) error {
	timer := d.String()
	newtimer := to.String()
	request := "UPDATE events SET timer = $1 WHERE timer = $2 AND information = $3"
	_, err := p.db.Exec(request, newtimer, timer, string(e))
	return err
}

func (p *dbProvder) FindEvents(parameters SearchParameters) ([]Event, error) {
	events := make([]Event, 0, 10)
	where := getWhereParameter(parameters)
	if where == "" {
		return events, fmt.Errorf("Invalid search parameters")
	}
	request := "SELECT information FROM events WHERE " + where
	rows, err := p.db.Query(request)
	if err != nil {
		return events, err
	}
	defer rows.Close()
	for rows.Next() {
		var info string
		err = rows.Scan(&info)
		if err != nil {
			return events, err
		}
		events = append(events, Event(info))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func getWhereParameter(p SearchParameters) string {
	if p.Year <= 0 {
		return ""
	}
	if p.Week > 0 {
		if p.Month == 0 && p.Day == 0 {
			return fmt.Sprintf("EXTRACT('week' from timer) = %d", p.Week)
		}
		return ""
	}
	if p.Month > 0 {
		if p.Day == 0 {
			return fmt.Sprintf("EXTRACT('month' from timer) = %d", p.Month)
		}
		if p.Day > 0 {
			return fmt.Sprintf("EXTRACT('month' from timer) = %d AND EXTRACT('day' from timer) = %d", p.Month, p.Day)
		}
	}
	return ""
}

func (p *dbProvder) Invoke(id string) {
	fmt.Println("Invoked!!!", id)
}
