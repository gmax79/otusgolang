package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gmax79/otusgolang/http/internal/calendar"
	"github.com/gmax79/otusgolang/http/internal/support"
	"go.uber.org/zap"
)

type httpCalendarAPI struct {
	server    *http.Server
	logger    *zap.Logger
	lasterror error
	calen     calendar.Calendar
}

func createServer(host string, zaplog *zap.Logger) *httpCalendarAPI {
	s := &httpCalendarAPI{logger: zaplog}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpRoot)
	mux.HandleFunc("/create_event", s.httpCreateEvent)
	mux.HandleFunc("/delete_event", s.httpDeleteEvent)
	//mux.HandleFunc("/update_event", s.httpUpdateEvent)
	mux.HandleFunc("/events_for_day", s.httpEventsForDay)
	mux.HandleFunc("/events_for_week", s.httpEventsForWeek)
	mux.HandleFunc("/events_for_month", s.httpEventsForMonth)
	s.server = &http.Server{Addr: host, Handler: mux}
	s.calen = calendar.Create()
	return s
}

func (s *httpCalendarAPI) logRequest(r *http.Request) {
	s.logger.Info("request", zap.String("method", r.Method), zap.String("url", r.URL.Path))
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

type dummyEvent struct {
	event string
}

func (d *dummyEvent) Invoke() {
	fmt.Println("Event !!!")
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
		if !events.AddEvent(&dummyEvent{event: event}) {
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
		if !s.calen.DeleteTrigger(time) {
			support.HTTPResponse(w, fmt.Errorf("event not found"))
			return
		}
		s.logger.Info("delete", zap.String("time", time))
		support.HTTPResponse(w, http.StatusOK)
		return
	}
	support.HTTPResponse(w, http.StatusBadRequest)
}

func (s *httpCalendarAPI) httpUpdateEvent(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	if pr := support.ReadPostRequest(r, w); pr != nil {

		//time := pr.Get("time")
		//event := pr.Get("event")
		//newtime := pr.Get("newtime")

		/*events, err := s.calen.AddTrigger(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		if !events.AddEvent(&dummyEvent{}) {
			support.HTTPResponse(w, http.StatusInternalServerError)
			return
		}*/
	}
	support.HTTPResponse(w, http.StatusBadRequest)
}

func (s *httpCalendarAPI) httpEventsForDay(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	support.HTTPResponse(w, http.StatusBadRequest)
}

func (s *httpCalendarAPI) httpEventsForWeek(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	support.HTTPResponse(w, http.StatusBadRequest)
}

func (s *httpCalendarAPI) httpEventsForMonth(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	support.HTTPResponse(w, http.StatusBadRequest)
}
