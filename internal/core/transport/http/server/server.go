package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"workout_app/docs/swagger"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type HTTPServer struct {
	mux    *http.ServeMux
	hanler http.Handler
	config Config
}

func NewHTTPServer(
	config Config,
	mux *http.ServeMux,
	middlewares ...func(http.Handler) http.Handler,
) *HTTPServer {
	handler := http.Handler(mux)

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return &HTTPServer{
		mux:    mux,
		hanler: handler,
		config: config,
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)

	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(swagger.SwaggerInfo.ReadDoc()))
		},
	)
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: s.hanler,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
	}

	return nil
}
