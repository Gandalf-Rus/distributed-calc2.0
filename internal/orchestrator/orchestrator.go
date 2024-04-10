package orchestrator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Gandalf-Rus/distributed-calc2.0/internal/config"
	l "github.com/Gandalf-Rus/distributed-calc2.0/internal/logger"
	"github.com/Gandalf-Rus/distributed-calc2.0/internal/middlewares"
	"github.com/gorilla/mux"
)

type Orchestrator struct {
	server *http.Server
}

func New() (*Orchestrator, error) {
	const (
		defaultHTTPServerWriteTimeout = time.Second * 15
		defaultHTTPServerReadTimeout  = time.Second * 15
	)

	orch := new(Orchestrator)

	router := mux.NewRouter()

	l.Logger.Info("initializing config...")
	err := config.InitConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	router.Use(middlewares.LoggingMiddleware)

	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hi")) }).Methods("GET")

	router.PathPrefix("/api").Handler(apiRouter)

	orch.server = &http.Server{
		Handler:      router,
		Addr:         ":" + strconv.Itoa(config.Cfg.ServerPort),
		WriteTimeout: defaultHTTPServerWriteTimeout,
		ReadTimeout:  defaultHTTPServerReadTimeout,
	}

	return orch, nil
}

func (o *Orchestrator) Run() error {
	l.Logger.Info("starting http server...")

	err := o.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server was stop with err: %w", err)
	}

	l.Logger.Info("server was stop")
	return nil
}

func (o *Orchestrator) stop(ctx context.Context) error {
	l.Logger.Info("shutdowning server...")
	err := o.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server was shutdown with error: %w", err)
	}
	l.Logger.Info("server was shutdown")
	return nil
}

func (o *Orchestrator) GracefulStop(serverCtx context.Context, sig <-chan os.Signal, serverStopCtx context.CancelFunc) {
	<-sig
	var timeOut = 30 * time.Second
	shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, timeOut)

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			l.Logger.Error("graceful shutdown timed out... forcing exit")
			os.Exit(1)
		}
	}()

	err := o.stop(shutdownCtx)
	if err != nil {
		l.Logger.Error("graceful shutdown timed out... forcing exit")
		os.Exit(1)
	}
	serverStopCtx()
	shutdownStopCtx()
}
