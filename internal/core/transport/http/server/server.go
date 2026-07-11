package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
}

func NewHTTPServer(
	config Config,
	mux *http.ServeMux,
) *HTTPServer {
	return &HTTPServer{
		mux:    mux,
		config: config,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: s.mux,
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
