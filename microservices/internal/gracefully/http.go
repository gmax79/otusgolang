package gracefully

import (
	"context"
	"net/http"
	"time"
)

// HTTPServer - wrapper to gracefully start and stop http server
type HTTPServer struct {
	server    *http.Server
	lastError error
}

// CreateHTTPServer - create server with host and main handler
func CreateHTTPServer(host string, handler http.Handler) *HTTPServer {
	var s HTTPServer
	s.server = &http.Server{Addr: host, Handler: handler}
	return &s
}

// Shutdown - command to stop server with timeout
func (s *HTTPServer) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// ListenAndServe - start listening
func (s *HTTPServer) ListenAndServe(timeToStart time.Duration) error {
	wait := make(chan struct{})
	go func() {
		close(wait)
		s.lastError = s.server.ListenAndServe()
	}()
	<-wait
	time.Sleep(timeToStart)
	return s.lastError
}

// GetLastError - return last error occurs in http server
func (s *HTTPServer) GetLastError() error {
	return s.lastError
}
