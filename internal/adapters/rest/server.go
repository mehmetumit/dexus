package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mehmetumit/dexus/internal/adapters/rest/middleware"
	"github.com/mehmetumit/dexus/internal/core/ports"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type ServerConfig struct {
	Addr string
	//Useful for maintaining connection between microservices
	//Stay open tls connection to decrease cost
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
type Server struct {
	config *ServerConfig
	logger ports.Logger
	app    ports.AppRunner
}

type Option func(http.Handler) http.Handler

func WithCacher(c ports.Cacher, l ports.Logger, ttl time.Duration) Option {
	cacheInterceptor := middleware.NewCacheInterceptor(c, l, ttl)
	return cacheInterceptor.InterceptHandler
}

func NewServer(logger ports.Logger, config *ServerConfig, app ports.AppRunner, options ...Option) {
	router := chi.NewRouter()
	router.Use(chiMiddleware.Logger, chiMiddleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// Load optional middlewares
	for _, op := range options {
		router.Use(op)
	}
	s := &http.Server{
		Addr:    config.Addr,
		Handler: router,
		// Useful for maintaining connection between microservices
		// Stay open tls connection to decrease cost
		IdleTimeout:  config.IdleTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	server := Server{
		config: config,
		logger: logger,
		app:    app,
	}
	server.initHandlers(router)
	go func() {
		logger.Debug("Starting server on", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			logger.Error(err)
			if err != http.ErrServerClosed {
				os.Exit(1)
			}
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	sig := <-sigChan
	logger.Info(fmt.Sprintf("SIGNAL:%s, Received terminate, graceful shutdown.", sig))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		//Release resources, close connections in here
		cancel() //  Releases resources if shutdown operation complete before timeout
	}()
	if err := s.Shutdown(ctx); err != nil {
		logger.Panic("Server shutdown failed:", err)
	}
}
func (s *Server) initHandlers(r *chi.Mux) {
	r.HandleFunc("/*", s.RedirectionHandler)
}
func (s *Server) RedirectionHandler(w http.ResponseWriter, r *http.Request) {
	to, err := s.app.FindRedirect(r.Context(), r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Don't cache redirections on client side by using 302 instead of 301
	http.Redirect(w, r, to.String(), http.StatusFound)
}
