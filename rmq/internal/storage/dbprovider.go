package storage

import (
	"database/sql"
	"fmt"

	"github.com/gmax79/otusgolang/rmq/internal/objects"
	"github.com/gmax79/otusgolang/rmq/internal/simple"
	_ "github.com/jackc/pgx/stdlib" // attach pgx postgres driver
)

// ConnectToDatabase - return standard sql.DB
func ConnectToDatabase(dsn string) (*sql.DB, error) {
	connection, err := sql.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, err
	}
	return connection, nil
}

// DbProvider main connection object
type DbProvider struct {
	db *sql.DB
}

// CreateProvider - return db provider by sql.DB connection
func CreateProvider(db *sql.DB) *DbProvider {
	return &DbProvider{db: db}
}

// GetTriggers - return triggers by date
func (p *DbProvider) GetTriggers() ([]simple.Date, error) {
	rows, err := p.db.Query("SELECT DISTINCT timer FROM events;")
	if err != nil {
		return []simple.Date{}, err
	}
	defer rows.Close()
	ids := make([]simple.Date, 0, 10)
	for rows.Next() {
		var timer string
		err := rows.Scan(&timer)
		if err != nil {
			return []simple.Date{}, err
		}
		var d simple.Date
		d.ParseDate(timer)
		ids = append(ids, d)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

// DeleteTrigger - delete trigger by date
func (p *DbProvider) DeleteTrigger(d simple.Date) error {
	timer := d.String()
	_, err := p.db.Exec("DELETE FROM events WHERE timer = $1;", timer)
	return err
}

// AddEvent - new event in db
func (p *DbProvider) AddEvent(d simple.Date, info string) error {
	timer := d.String()
	_, err := p.db.Exec("INSERT INTO events (timer, information) VALUES($1, $2) ON CONFLICT (timer, information) DO NOTHING;", timer, info)
	return err
}

// GetEventsCount - count events by date
func (p *DbProvider) GetEventsCount(d simple.Date) (int, error) {
	timer := d.String()
	var count int
	err := p.db.QueryRow("SELECT COUNT (*) FROM events WHERE timer = $1", timer).Scan(&count)
	return count, err
}

// DeleteEventIndex - delete event by index
func (p *DbProvider) DeleteEventIndex(d simple.Date, index int) error {
	request := "DELETE FROM events WHERE ctid IN (SELECT ctid FROM events WHERE timer = $1::timestamp limit 1 offset $2);"
	timer := d.String()
	_, err := p.db.Exec(request, timer, index)
	return err
}

// DeleteEvent - delete event by date
func (p *DbProvider) DeleteEvent(d simple.Date, e objects.Event) error {
	timer := d.String()
	_, err := p.db.Exec("DELETE FROM events WHERE timer = $1::timestamp AND information = $2;", timer, string(e))
	return err
}

// GetEvent - return event by date and index
func (p *DbProvider) GetEvent(d simple.Date, index int) (objects.Event, error) {
	request := "SELECT information FROM events WHERE timer = $1::timestamp limit 1 offset $2;"
	timer := d.String()
	var e objects.Event
	err := p.db.QueryRow(request, timer, index).Scan(&e)
	return e, err
}

// MoveEvent - move event from date to another date
func (p *DbProvider) MoveEvent(d simple.Date, e objects.Event, to simple.Date) error {
	timer := d.String()
	newtimer := to.String()
	request := "UPDATE events SET timer = $1 WHERE timer = $2 AND information = $3"
	_, err := p.db.Exec(request, newtimer, timer, string(e))
	return err
}

// FindEvents - find events by search parameters
func (p *DbProvider) FindEvents(parameters objects.SearchParameters) ([]objects.Event, error) {
	events := make([]objects.Event, 0, 10)
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
		events = append(events, objects.Event(info))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func getWhereParameter(p objects.SearchParameters) string {
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

// Invoke - event happend method
func (p *DbProvider) Invoke(id string) {
	fmt.Println("Invoked!!!", id)
}
