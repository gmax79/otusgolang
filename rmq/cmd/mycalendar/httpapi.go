package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gmax79/otusgolang/rmq/internal/calendar"
	"github.com/gmax79/otusgolang/rmq/internal/support"
	"go.uber.org/zap"
)

type httpCalendarAPI struct {
	server    *http.Server
	logger    *zap.Logger
	lastError error
	calen     calendar.Calendar
}

func createServer(calen calendar.Calendar, host string, zaplog *zap.Logger) *httpCalendarAPI {
	s := &httpCalendarAPI{logger: zaplog}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpRoot)
	mux.HandleFunc("/create_event", s.httpCreateEvent)
	mux.HandleFunc("/delete_event", s.httpDeleteEvent)
	mux.HandleFunc("/move_event", s.httpMoveEvent)
	mux.HandleFunc("/events_for_day", s.httpEventsForDay)
	mux.HandleFunc("/events_for_week", s.httpEventsForWeek)
	mux.HandleFunc("/events_for_month", s.httpEventsForMonth)
	s.server = &http.Server{Addr: host, Handler: mux}
	s.calen = calen
	return s
}

func (s *httpCalendarAPI) logRequest(r *http.Request) {
	url := r.URL.Path
	if r.Method == http.MethodGet {
		url = r.URL.RequestURI()
	}
	s.logger.Info("request", zap.String("method", r.Method), zap.String("url", url))
}

func (s *httpCalendarAPI) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)
}

func (s *httpCalendarAPI) ListenAndServe() {
	go func() {
		s.lastError = s.server.ListenAndServe()
	}()
}

func (s *httpCalendarAPI) GetLastError() error {
	return s.lastError
}

func (s *httpCalendarAPI) httpRoot(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	support.HTTPResponse(w, http.StatusNotFound)
}

func (s *httpCalendarAPI) httpCreateEvent(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	if pr := support.ReadPostRequest(r, w); pr != nil {
		time := pr.Get("time")
		event := pr.Get("event")
		if time == "" || event == "" {
			support.HTTPResponse(w, http.StatusBadRequest)
			return
		}
		t, err := calendar.ParseValidDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		events, err := s.calen.AddTrigger(t)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		err = events.AddEvent(calendar.Event(event))
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		s.logger.Info("new", zap.String("time", t.String()), zap.String("event", event))
		support.HTTPResponse(w, http.StatusOK)
		return
	}
	support.HTTPResponse(w, http.StatusBadRequest)
}

func (s *httpCalendarAPI) httpDeleteEvent(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	if pr := support.ReadPostRequest(r, w); pr != nil {
		time := pr.Get("time")
		event := pr.Get("event")
		if time == "" || event == "" {
			support.HTTPResponse(w, http.StatusBadRequest)
			return
		}
		t, err := calendar.ParseValidDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		events, err := s.calen.GetEvents(t)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		err = events.DeleteEvent(calendar.Event(event))
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		s.logger.Info("delete", zap.String("time", time), zap.String("event", event))
		support.HTTPResponse(w, http.StatusOK)
		return
	}
	support.HTTPResponse(w, http.StatusBadRequest)
}

func (s *httpCalendarAPI) httpMoveEvent(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	if pr := support.ReadPostRequest(r, w); pr != nil {
		time := pr.Get("time")
		event := pr.Get("event")
		newtime := pr.Get("newtime")
		if time == "" || event == "" || newtime == "" {
			support.HTTPResponse(w, http.StatusBadRequest)
			return
		}
		t, err := calendar.ParseValidDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		newt, err := calendar.ParseValidDate(newtime)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		events, err := s.calen.GetEvents(t)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		err = events.MoveEvent(calendar.Event(event), newt)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		support.HTTPResponse(w, http.StatusOK)
		return
	}
	support.HTTPResponse(w, http.StatusBadRequest)
}

func (s *httpCalendarAPI) httpEventsForDay(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	if pr := support.ReadGetRequest(r, w); pr != nil {
		time := pr.Get("day")
		if time == "" {
			support.HTTPResponse(w, http.StatusBadRequest)
			return
		}
		d, err := calendar.ParseValidDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		var sp calendar.SearchParameters
		sp.Day = d.Day
		sp.Month = d.Month
		sp.Year = d.Year

		count, err := s.calen.FindEvents(sp)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		result := fmt.Sprintf("{ \"result\": %d }", len(count))
		support.HTTPResponse(w, result)
		return
	}
	support.HTTPResponse(w, http.StatusBadRequest)
}

func (s *httpCalendarAPI) httpEventsForWeek(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	if pr := support.ReadGetRequest(r, w); pr != nil {
		time := pr.Get("week")
		if time == "" {
			support.HTTPResponse(w, http.StatusBadRequest)
			return
		}
		d, err := calendar.ParseDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		var sp calendar.SearchParameters
		sp.Week = d.Month
		sp.Year = d.Year

		count, err := s.calen.FindEvents(sp)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		result := fmt.Sprintf("{ \"result\": %d }", len(count))
		support.HTTPResponse(w, result)
		return
	}
	support.HTTPResponse(w, http.StatusBadRequest)
}

func (s *httpCalendarAPI) httpEventsForMonth(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	if pr := support.ReadGetRequest(r, w); pr != nil {
		time := pr.Get("month")
		if time == "" {
			support.HTTPResponse(w, http.StatusBadRequest)
			return
		}
		d, err := calendar.ParseDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		var sp calendar.SearchParameters
		sp.Month = d.Month
		sp.Year = d.Year

		count, err := s.calen.FindEvents(sp)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		result := fmt.Sprintf("{ \"result\": %d }", len(count))
		support.HTTPResponse(w, result)
		return
	}
	support.HTTPResponse(w, http.StatusBadRequest)
}
