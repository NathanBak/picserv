package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// A Server is a web server that extends http.Server.  It should be created with the New() function,
// started with the Server.ListenAndServe() function, and can be cleanly closed with the
// Server.Shutdown() function.
type Server struct {
	http.Server
	log Logger
	cfg Config
}

// New creates, configures, and returns a new server instance.
func New(cfg Config) (*Server, error) {

	if cfg.Logger == nil {
		cfg.Logger = defaultLogger{}
	}

	flip, err := newFlipper(cfg.Picker, cfg.PicLife, cfg.Logger)
	if err != nil {
		return nil, err
	}
	cfg.Picker = flip

	s := &Server{
		Server: http.Server{
			Addr:         fmt.Sprintf(":%v", cfg.Port),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		cfg: cfg,
		log: cfg.Logger,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.
		Methods(http.MethodGet).
		Path("/next").
		Name("Next Pic").
		HandlerFunc(s.next)
	s.Server.Handler = router

	return s, nil
}

func (s *Server) ListenAndServe() error {
	msg := fmt.Sprintf("Starting server on port %d", s.cfg.Port)
	s.log.Info(msg)
	return s.Server.ListenAndServe()
}

// ServeHttp routes and handles a request.  It is typically used for testing purposes.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}

// Shutdown cleanly shutsdown the server and blocks until complete (or timing out).
func (s *Server) Shutdown(ctx context.Context) error {
	msg := "Shutting down server"
	s.log.Info(msg)
	return s.Server.Shutdown(ctx)
}
