package middleware

import (
	"context"
	"net/http"
	"time"
)

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		go func() {
			next.ServeHTTP(w, r)
			close(done)
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			w.WriteHeader(500)
			w.Write([]byte(`{"message": "request timed out"}`))
		}
	})
}
