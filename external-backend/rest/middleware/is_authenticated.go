package middleware

import "net/http"

func IsAuthenticated() func(next http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return authenticator{
			handler: handler,
		}
	}
}

type authenticator struct {
	handler http.Handler
}

func (a authenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
