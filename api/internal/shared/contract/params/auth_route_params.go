package params

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"upworkapi/pkg/supa"
)

type AuthRouteParams struct {
	RootGroup     *echo.Group
	SessionSecret string
	Logger        *slog.Logger
	Supabase      *supa.Supabase
}
