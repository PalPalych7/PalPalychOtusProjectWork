package internalhttp

import (
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		myTimeStart := time.Now()
		next.ServeHTTP(rw, r)
		s.myLogger.Info(r.RemoteAddr, myTimeStart, r.Method, r.Proto, time.Since(myTimeStart), r.UserAgent())
	})
}
