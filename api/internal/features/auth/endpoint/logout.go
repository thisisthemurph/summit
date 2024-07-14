package endpoint

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"upworkapi/internal/shared/auth"
	"upworkapi/internal/shared/contract"
	"upworkapi/internal/shared/contract/params"
)

type logOutEndpoint struct {
	params.AuthRouteParams
}

func NewLogOutEndpoint(routeParams params.AuthRouteParams) contract.Endpoint {
	return &logOutEndpoint{
		AuthRouteParams: routeParams,
	}
}

func (ep *logOutEndpoint) MapEndpoint() {
	ep.RootGroup.POST("/logout", ep.handleLogOut())
}

func (ep *logOutEndpoint) handleLogOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		store, _ := auth.GetCookieStoreSession(c.Request(), auth.SessionCookieStoreKey, ep.SessionSecret)
		store.Values[auth.AccessTokenKey] = ""
		store.Options.MaxAge = -1
		if err := store.Save(c.Request(), c.Response()); err != nil {
			ep.Logger.Error("logout failed. could not delete cookie store", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
