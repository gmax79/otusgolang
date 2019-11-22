package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gmax79/otusgolang/http/internal/calendar"
	"go.uber.org/zap"
)

// Response - send http answer to client
func Response(w http.ResponseWriter, v interface{}) {
	switch v.(type) {
	case nil:
		w.WriteHeader(http.StatusNotFound)
	case int:
		w.WriteHeader(v.(int))
	case error:
		err := v.(error)
		text := "{ \"error\" : " + err.Error() + "  }"
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(text))
	default:
		answer, err := json.Marshal(v)
		if err != nil {
			Response(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(answer)
		}
	}
}

type httpCalandarAPI struct {
	server    *http.Server
	logger    *zap.Logger
	lasterror error
	clr       calendar.Calendar
}

type createEventRequest struct {
	ID string
}

func (p *createEventRequest) ReadParameters(r *http.Request, w http.ResponseWriter) (ok bool) {
	if r.Method != http.MethodPost {
		Response(w, http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		Response(w, err)
		return
	}
	p.ID = r.Form.Get("id")
	return true
}

func createServer(host string, zaplog *zap.Logger) *httpCalandarAPI {
	s := &httpCalandarAPI{logger: zaplog}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpRoot)
	mux.HandleFunc("/create_event", s.httpCreateEvent)
	s.server = &http.Server{Addr: host, Handler: mux}
	s.clr = calendar.CreateCalendar()
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
	Response(w, nil)
}

func (s *httpCalandarAPI) httpCreateEvent(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	var p createEventRequest
	if !p.ReadParameters(r, w) {
		return
	}
	fmt.Fprint(w, "id = ", p.ID)
	Response(w, http.StatusOK)
}
