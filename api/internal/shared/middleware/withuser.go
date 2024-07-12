package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"upworkapi/cmd/api/application"
	"upworkapi/internal/shared/auth"
	"upworkapi/pkg/supa"
)

type AuthMiddleware struct {
	Config   *application.Config
	Supabase *supa.Supabase
	Logger   *slog.Logger
}

func NewAuthMiddleware(config *application.Config, sb *supa.Supabase, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		Config:   config,
		Supabase: sb,
		Logger:   logger,
	}
}

// WithUser middleware sets an AuthenticatedUser on the context if a session auth is present.
func (m AuthMiddleware) WithUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		session, err := auth.GetCookieStoreSession(c.Request(), auth.SessionCookieStoreKey, m.Config.SessionSecret)
		if err != nil {
			return next(c)
		}

		sessionAccessToken := session.Values[auth.AccessTokenKey]
		accessToken, ok := sessionAccessToken.(string)
		if !ok {
			return next(c)
		}

		sbUser, err := m.Supabase.Client.Auth.User(ctx, accessToken)
		if err != nil {
			m.Logger.Info("supabase error authenticating user", "error", err)
			return next(c)
		}

		userID, err := uuid.Parse(sbUser.ID)
		if err != nil {
			m.Logger.Info("supabase user ID is not a valid UUID", "error", err)
			return next(c)
		}

		authenticatedAccount := auth.AuthenticatedUser{
			ID:          userID,
			AccessToken: accessToken,
			LoggedIn:    true,
		}

		auth.SaveAuthenticatedUserInContext(c, authenticatedAccount)
		return next(c)
	}
}

func (m AuthMiddleware) WithAuthenticatedUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := auth.GetAuthenticatedUser(c)
		if !user.LoggedIn {
			m.Logger.Info("User not authenticated, redirecting")
			return c.JSON(http.StatusForbidden, echo.Map{"error": "forbidden"})
		}
		return next(c)
	}
}
