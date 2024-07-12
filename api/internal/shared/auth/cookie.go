package auth

import (
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	SessionCookieStoreKey = "cs"
	AccessTokenKey        = "accessToken"
)

type ContextKey string

const AuthenticatedUserContextKey ContextKey = "authenticatedUser"

type AuthenticatedUser struct {
	ID          uuid.UUID `json:"id"`
	AccessToken string    `json:"access_token"`
	LoggedIn    bool      `json:"logged_in"`
}

func SetAuthSession(c echo.Context, accessToken, sessionSecret string) error {
	session, _ := GetCookieStoreSession(c.Request(), SessionCookieStoreKey, sessionSecret)
	session.Values[AccessTokenKey] = accessToken
	session.Options.HttpOnly = true
	return session.Save(c.Request(), c.Response().Writer)
}

func GetCookieStoreSession(r *http.Request, cookieStoreName, secret string) (*sessions.Session, error) {
	store := sessions.NewCookieStore([]byte(secret))
	return store.Get(r, cookieStoreName)
}
