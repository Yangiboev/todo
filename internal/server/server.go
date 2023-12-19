package server

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/Yangiboev/todo/config"
	"github.com/Yangiboev/todo/pkg/logger"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server struct
type Server struct {
	echo   *echo.Echo
	cfg    *config.Config
	db     *sqlx.DB
	logger logger.Logger
}

// NewServer constructor
func NewServer(cfg *config.Config, db *sqlx.DB, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, db: db, logger: logger}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	go func() {
		s.logger.Infof("Starting Debug Server on PORT: %s", s.cfg.Server.PprofPort)
		if err := http.ListenAndServe(s.cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
