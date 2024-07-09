package params

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"upworkapi/pkg/supa"
)

type AuthRouteParams struct {
	LoginGroup  *echo.Group
	SignupGroup *echo.Group
	LogoutGroup *echo.Group

	Logger   *slog.Logger
	Supabase *supa.Supabase
}
