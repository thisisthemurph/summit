package endpoint

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
	"net/http"
	"upworkapi/internal/features/auth/cookie"
	"upworkapi/internal/shared/contract"
	"upworkapi/internal/shared/contract/params"
)

var (
	ErrLoggingIn          = errors.New("there has been an issue logging you in")
	ErrInvalidCredentials = errors.New("the provided credentials are invalid")
)

type loginEndpoint struct {
	params.AuthRouteParams
}

func NewLoginEndpoint(routeParams params.AuthRouteParams) contract.Endpoint {
	return &loginEndpoint{
		routeParams,
	}
}

func (ep *loginEndpoint) MapEndpoint() {
	ep.LoginGroup.POST("", ep.loginHandler())
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (ep *loginEndpoint) loginHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		authDetails, err := ep.Supabase.Client.Auth.SignIn(ctx, supabase.UserCredentials{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			ep.Logger.Warn("login failed", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, ErrInvalidCredentials.Error())
		}

		if err := cookie.SetAuthSession(c, authDetails.AccessToken, "session-secret"); err != nil {
			ep.Logger.Error("could not set auth session", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, ErrLoggingIn.Error())
		}

		return c.JSON(http.StatusOK, echo.Map{})
	}
}
