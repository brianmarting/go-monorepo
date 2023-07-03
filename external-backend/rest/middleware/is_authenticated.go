package middleware

import (
	"fmt"
	"net/http"
	"os"
)

var (
	ipsHost = os.Getenv("IPS_HOST")
	ipsPort = os.Getenv("IPS_PORT")
)

func IsAuthenticated() func(next http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return authenticator{
			handler: handler,
			client:  http.Client{},
		}
	}
}

type authenticator struct {
	handler http.Handler
	client  http.Client
}

func (a authenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("authorization")

	if authorization == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s:%s", ipsHost, ipsPort), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	request.Header.Add("authorization", authorization)

	response, err := a.client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if response.StatusCode != 200 {
		http.Error(w, fmt.Sprintf("unauthorized with code %d", response.StatusCode), http.StatusUnauthorized)
		return
	}

	a.handler.ServeHTTP(w, r)
}
