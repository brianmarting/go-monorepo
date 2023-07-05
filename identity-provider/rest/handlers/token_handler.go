package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-monorepo/identity-provider/service"
	"go-monorepo/internal/model"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessTokenSecret  = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
	refreshToken       = "rt"
)

type TokenHandler interface {
	Validate() http.HandlerFunc
	PostRefreshToken() http.HandlerFunc
}

type tokenHandler struct {
	userService service.UserService
}

func NewTokenHandler() TokenHandler {
	return &tokenHandler{
		userService: service.NewUserService(),
	}
}

func (t tokenHandler) Validate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(
				w,
				errors.New("no token is present").Error(),
				http.StatusBadRequest,
			)
			return
		}

		rawToken := strings.Split(token, " ")
		if len(rawToken) < 2 {
			http.Error(
				w,
				errors.New("invalid token").Error(),
				http.StatusBadRequest,
			)
			return
		}

		if _, err := parseToken(rawToken[1], accessTokenSecret); err != nil {
			http.Error(
				w,
				fmt.Errorf("invalid token: %v", err).Error(),
				http.StatusBadRequest,
			)
			return
		}

		w.WriteHeader(200)
	}
}

func (t tokenHandler) PostRefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := retrieveClaimsFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, tokenVersion := claims["id"], claims["tokenVersion"]
		if id == "" || tokenVersion == "" {
			http.Error(w, errors.New("failed to retrieve data from claims").Error(), http.StatusBadRequest)
			return
		}
		user, err := t.userService.GetByExternalId(id.(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.TokenVersion != tokenVersion {
			http.Error(w, errors.New("invalid token version").Error(), http.StatusBadRequest)
			return
		}

		accessToken, err := CreateAccessToken(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.SetCookie(w, CreateCookie(user))
		if _, err = w.Write(accessToken); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func CreateCookie(user model.User) *http.Cookie {
	signedToken, _ := signRefreshToken(user)
	return &http.Cookie{
		Name:     refreshToken,
		Value:    signedToken,
		Path:     "/",
		HttpOnly: false,
	}
}

func CreateAccessToken(user model.User) ([]byte, error) {
	signedAccessToken, err := signAccessToken(user)
	if err != nil {
		return nil, err
	}

	return json.Marshal(model.AccessToken{
		AccessToken: fmt.Sprintf("Bearer %s", signedAccessToken),
	})
}

func signAccessToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ExternalId,
		"name": user.Name,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	})
	return token.SignedString(accessTokenSecret)
}

func signRefreshToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":           user.ExternalId,
		"tokenVersion": user.TokenVersion,
		"exp":          time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(refreshTokenSecret)
}

func retrieveClaimsFromRequest(r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie(refreshToken)
	if err != nil {
		return nil, err
	}

	token, err := parseToken(cookie.String(), refreshTokenSecret)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to retrieve data from claims")
	}
	return claims, nil
}

func parseToken(tokenString string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("failed to validate token")
		}

		return secret, nil
	})
}
