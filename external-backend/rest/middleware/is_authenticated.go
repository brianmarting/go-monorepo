package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
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
	authorization := r.Header.Get("Authorization")

	if authorization == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://%s:%s/user/token/validate", ipsHost, ipsPort), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	request.Header.Add("Authorization", authorization)

	response, err := a.client.Do(request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "timed out when validating token", http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if response.StatusCode != 200 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	a.handler.ServeHTTP(w, r)
}
