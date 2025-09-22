package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/buni/tx-parser/internal/pkg/configuration"
	"github.com/buni/tx-parser/internal/pkg/sloglog"
)

// Server ...
type Server struct {
	Context     context.Context
	cancel      func()
	Logger      *slog.Logger
	Config      *configuration.Configuration
	Router      chi.Router
	httpServer  *http.Server
	done        chan os.Signal
	gracePeriod time.Duration
}

// Option ...
type Option func(*Server) error

// WithGracePeriod ...
func WithGracePeriod(gracePeriod time.Duration) Option {
	return func(s *Server) error {
		s.gracePeriod = gracePeriod

		return nil
	}
}

// NewServer ...
func NewServer(ctx context.Context, opts ...Option) (server *Server, err error) {
	server = &Server{
		gracePeriod: 1 * time.Second,
		done:        make(chan os.Signal, 1),
	}

	signal.Notify(server.done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	server.Config, err = configuration.NewConfiguration()
	if err != nil {
		return nil, fmt.Errorf("create configuration: %w", err)
	}

	server.Context, server.cancel = context.WithCancel(ctx)

	zapLogger, err := zap.NewProduction(zap.WithCaller(true), zap.AddStacktrace(zapcore.DPanicLevel))
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	server.Logger = slog.New(sloglog.ApplyMiddleware(zapslog.NewHandler(zapLogger.Core(), zapslog.WithCaller(true))))

	slog.SetDefault(server.Logger)

	zap.ReplaceGlobals(zapLogger)

	server.Context = sloglog.ToContext(server.Context, server.Logger)

	server.Router = chi.NewRouter()

	server.Router.Use(middleware.Recoverer)
	server.Router.Use(sloglog.RequestLogger())

	for _, opt := range opts {
		err = opt(server)
		if err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	return server, nil
}

// Start ...
func (a *Server) Start() error {
	a.Logger.InfoContext(a.Context, "starting server")
	a.httpServer = &http.Server{
		Addr:              a.Config.ToHost(),
		Handler:           h2c.NewHandler(a.Router, &http2.Server{}),
		ReadHeaderTimeout: 120 * time.Second,
		IdleTimeout:       120 * time.Second,
		BaseContext:       func(_ net.Listener) context.Context { return a.Context },
	}

	go func() {
		err := a.httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			a.Logger.ErrorContext(a.Context, "http server error", sloglog.Error(err))
		}
	}()

	return nil
}

// Shutdown ...
func (a *Server) Wait(shutdownFuncs ...func()) {
	<-a.done
	ctx, cancel := context.WithTimeout(a.Context, a.gracePeriod)
	defer cancel()

	a.Logger.InfoContext(a.Context, "shutting down")
	a.httpServer.Shutdown(ctx) //nolint:errcheck,revive
	for _, shutdownFunc := range shutdownFuncs {
		shutdownFunc()
	}

	a.cancel()
	<-ctx.Done()
	a.Logger.InfoContext(a.Context, "shutdown")
}
