package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gmax79/otusgolang/microservices/internal/gracefully"
	"github.com/gmax79/otusgolang/microservices/internal/pmetrics"

	"github.com/gmax79/otusgolang/microservices/internal/calendar"
	"github.com/gmax79/otusgolang/microservices/internal/objects"
	"github.com/gmax79/otusgolang/microservices/internal/simple"
	"github.com/gmax79/otusgolang/microservices/internal/support"

	"go.uber.org/zap"
)

type httpCalendarAPI struct {
	server *gracefully.HTTPServer
	logger *zap.Logger
	calen  calendar.Calendar
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
	metricsHandler := pmetrics.AttachPrometheusToHandler(mux)
	s.server = gracefully.CreateHTTPServer(host, metricsHandler)
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
	err := s.server.Shutdown(3 * time.Second)
	if err != nil {
		s.logger.Error("shutdown", zap.String("error", err.Error()))
	}
}

func (s *httpCalendarAPI) ListenAndServe() error {
	timetostart := time.Millisecond * 200
	return s.server.ListenAndServe(timetostart)
}

func (s *httpCalendarAPI) GetLastError() error {
	return s.server.GetLastError()
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
		t, err := simple.ParseValidDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		events, err := s.calen.AddTrigger(t)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}

		var newevent objects.Event
		newevent.Alerttime = t
		newevent.Information = event

		err = events.AddEvent(newevent)
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
		t, err := simple.ParseValidDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		events, err := s.calen.GetEvents(t)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}

		var delevent objects.Event
		delevent.Alerttime = t
		delevent.Information = event
		err = events.DeleteEvent(delevent)
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
		t, err := simple.ParseValidDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		newt, err := simple.ParseValidDate(newtime)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		events, err := s.calen.GetEvents(t)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}

		var moveevent objects.Event
		moveevent.Alerttime = t
		moveevent.Information = event
		err = events.MoveEvent(moveevent, newt)
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
		d, err := simple.ParseValidDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		var sp objects.SearchParameters
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
		d, err := simple.ParseDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		var sp objects.SearchParameters
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
		d, err := simple.ParseDate(time)
		if err != nil {
			support.HTTPResponse(w, err)
			return
		}
		var sp objects.SearchParameters
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
