package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gmax79/otusgolang/http/internal/calendar"
	"github.com/gmax79/otusgolang/http/internal/support"
	"go.uber.org/zap"
)

type httpCalandarAPI struct {
	server    *http.Server
	logger    *zap.Logger
	lasterror error
	cr        calendar.Calendar
}

func createServer(host string, zaplog *zap.Logger) *httpCalandarAPI {
	s := &httpCalandarAPI{logger: zaplog}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpRoot)
	mux.HandleFunc("/create_event", s.httpCreateEvent)
	s.server = &http.Server{Addr: host, Handler: mux}
	s.cr = calendar.CreateCalendar()
	return s
}

func (s *httpCalandarAPI) logRequest(r *http.Request) {
	s.logger.Info("request", zap.String("url", r.URL.Path))
}

func (s *httpCalandarAPI) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)
}

func (s *httpCalandarAPI) ListenAndServe() {
	go func() {
		s.lasterror = s.server.ListenAndServe()
	}()
}

func (s *httpCalandarAPI) GetLastError() error {
	return s.lasterror
}

func (s *httpCalandarAPI) httpRoot(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	support.HTTPResponse(w, http.StatusNotFound)
}

func (s *httpCalandarAPI) httpCreateEvent(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	if pr := support.ReadPostRequest(r, w); r != nil {
		id := pr.Get("id")
		_ = id

		t := calendar.CreateCalendarTrigger()
		//t.AddTrigger
	}
}

func (s *httpCalandarAPI) httpDeleteEvent(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)

}
