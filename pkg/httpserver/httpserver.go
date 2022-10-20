package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"
)

type Handler interface {
	Handler(builder middleware.Builder) http.Handler
}

type Server struct {
	// routes store
	mp map[string]Handler

	// global middleware
	builder middleware.Builder
}

func New() *Server {
	return &Server{
		mp: map[string]Handler{},
	}
}

func (s *Server) AddWithFullPath(path string, handler Handler) {
	s.mp[path+"/"] = handler
}

func (s *Server) AddWithCurrentPath(path string, handler Handler) {
	s.mp[path] = handler
}

func (s *Server) Middleware(builder middleware.Builder) {
	s.builder = builder
}

func (s *Server) Serve(addr string) error {
	for p, h := range s.mp {
		http.Handle(p, metric(h.Handler(nil)))
	}

	return http.ListenAndServe(addr, nil)
}

func metric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != any(nil) {
				panicfields := logrus.WithFields(logrus.Fields{
					"method": r.Method,
					"url":    r.RequestURI,
					"panic":  true,
					"msg":    fmt.Sprintf("%+v", err),
				})

				panicfields.Info()
			}
		}()

		start := time.Now()
		next.ServeHTTP(w, r)

		fields := logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.RequestURI,
			"time":   time.Since(start),
		})

		fields.Info()
	})
}
