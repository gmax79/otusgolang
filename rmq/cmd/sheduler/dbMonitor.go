package main

import (
	"database/sql"
	"time"

	"github.com/gmax79/otusgolang/rmq/internal/simple"
	_ "github.com/jackc/pgx/stdlib" // attach pgx postgres driver
)

func connectToDatabase(dsn string, finished chan<- string) (*dbMonitor, error) {
	connection, err := sql.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, err
	}
	timers := make(map[time.Time]string)
	return &dbMonitor{db: connection, timers: timers, finished: finished}, nil
}

type dbMonitor struct {
	db       *sql.DB
	timers   map[time.Time]string
	finished chan<- string
}

func (m *dbMonitor) Close() {
	m.db.Close()
}

func (m *dbMonitor) ReadEvents() error {
	rows, err := m.db.Query("SELECT timer, information FROM events")
	if err != nil {
		return err
	}
	defer rows.Close()

	newtimers := map[time.Time]string{}
	now := simple.NowDate()
	for rows.Next() {
		var timer time.Time
		var info string
		if err = rows.Scan(&timer, &info); err != nil {
			return err
		}
		if now.Before(timer) {
			newtimers[timer] = info
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}
	for nt, e := range newtimers {
		if _, ok := m.timers[nt]; !ok {
			go func(event string, duration time.Duration) {
				t := time.NewTimer(duration)
				<-t.C
				m.finished <- event
			}(e, nt.Sub(now))
		}
	}
	m.timers = newtimers
	return nil
}

func (m *dbMonitor) GetNearestEvent() (event time.Time, ok bool) {
	if len(m.timers) == 0 {
		return
	}
	for t := range m.timers {
		event = t
		break
	}
	for t := range m.timers {
		if t.Before(event) {
			event = t
		}
	}
	return event, true
}
