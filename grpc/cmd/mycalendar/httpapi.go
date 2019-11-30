package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gmax79/otusgolang/grpc/internal/calendar"
	"github.com/gmax79/otusgolang/grpc/internal/support"
	"go.uber.org/zap"
)

type httpCalendarAPI struct {
	server    *http.Server
	logger    *zap.Logger
	lasterror error
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
		s.lasterror = s.server.ListenAndServe()
	}()
}

func (s *httpCalendarAPI) GetLastError() error {
	return s.lasterror
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
		events, err := s.calen.AddTrigger(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		if !events.AddEvent(&SimpleEvent{event: event}) {
			support.HTTPResponse(w, http.StatusBadRequest)
			return
		}
		tm, _ := s.calen.GetTriggerAlert(time)
		s.logger.Info("new", zap.String("time", tm.String()), zap.String("event", event))
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
		events := s.calen.GetEvents(time)
		if events == nil {
			support.HTTPResponse(w, fmt.Errorf("trigger %s not found", time))
			return
		}
		index := events.FindEvent(event)
		if index == -1 {
			support.HTTPResponse(w, fmt.Errorf("event %s in trigger %s not found", event, time))
			return
		}
		events.DeleteEvent(index)
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
		events := s.calen.GetEvents(time)
		if events == nil {
			support.HTTPResponse(w, fmt.Errorf("trigger %s not found", time))
			return
		}
		index := events.FindEvent(event)
		if index == -1 {
			support.HTTPResponse(w, fmt.Errorf("event %s in trigger %s not found", event, time))
			return
		}
		newevents, err := s.calen.AddTrigger(newtime)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		foundevent := events.GetEvent(index)
		events.DeleteEvent(index)
		newevents.AddEvent(foundevent)
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
		var d calendar.Date
		err := d.ParseDate(time)
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
		result := fmt.Sprintf(`{ "result": %d }`, len(count))
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
		var d calendar.Date
		err := d.ParseDate(time)
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
		result := fmt.Sprintf(`{ "result": %d }`, len(count))
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
		var d calendar.Date
		err := d.ParseDate(time)
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
		result := fmt.Sprintf(`{ "result": %d }`, len(count))
		support.HTTPResponse(w, result)
		return
	}
	support.HTTPResponse(w, http.StatusBadRequest)
}
