package params

import (
	"github.com/labstack/echo/v4"
	"log/slog"
)

type OnboardingRouteParams struct {
	OnboardingGroup *echo.Group
	Logger          *slog.Logger
}
