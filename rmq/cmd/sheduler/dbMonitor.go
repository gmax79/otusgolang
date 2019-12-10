package main

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/stdlib" // attach pgx postgres driver
)

func connectToDatabase(dsn string) (*dbMonitor, error) {
	connection, err := sql.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, err
	}
	timers := make(map[time.Time]string)
	return &dbMonitor{db: connection, timers: timers}, nil
}

type dbMonitor struct {
	db     *sql.DB
	timers map[time.Time]string
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
	now := time.Now()
	for rows.Next() {
		var timer time.Time
		var info string
		if err = rows.Scan(&timer, &info); err != nil {
			return err
		}
		if now.After(timer) {
			m.timers[timer] = info
		}
	}
	return nil
}

func (m *dbMonitor) SelectNextEvent() {
}
