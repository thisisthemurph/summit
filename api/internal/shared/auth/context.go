package auth

import (
	"context"
	"github.com/labstack/echo/v4"
)

func GetAuthenticatedUser(c echo.Context) AuthenticatedUser {
	user, ok := c.Get(string(AuthenticatedUserContextKey)).(AuthenticatedUser)
	if !ok {
		return AuthenticatedUser{}
	}
	return user
}

func SaveAuthenticatedUserInContext(c echo.Context, user AuthenticatedUser) {
	r := c.Request()
	ctx := c.Request().Context()

	// Persist session in the echo context
	c.Set(string(AuthenticatedUserContextKey), user)
	// Persist the session in the request context
	c.SetRequest(r.WithContext(context.WithValue(ctx, AuthenticatedUserContextKey, user)))
}
