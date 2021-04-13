package internalhttp

import (
	"context"
	"net"
	"net/http"

	"github.com/FedoseevAlex/OtusGolangHomeWork/hw12_13_14_15_calendar/internal/app"
)

type Server struct { // TODO
	httpServer http.Server
	log        app.Logger
}

type Application interface { // TODO
	Logger() app.Logger
}

func versionHandler(app Application) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, err := rw.Write([]byte("Hi, this is calendar app!"))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func NewServer(app Application, host, port string) *Server {
	router := http.NewServeMux()
	router.Handle("/hello", loggingMiddleware(versionHandler(app), app.Logger()))
	return &Server{
		httpServer: http.Server{
			Addr:    net.JoinHostPort(host, port),
			Handler: router,
		},
		log: app.Logger(),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Debug(
		"HTTP server started",
		map[string]interface{}{
			"address": s.httpServer.Addr,
		})

	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Debug("HTTP server stopped")
	return s.httpServer.Shutdown(ctx)
}
