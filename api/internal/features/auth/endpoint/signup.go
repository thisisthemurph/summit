package endpoint

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/nedpals/supabase-go"
	"net/http"
	"upworkapi/internal/features/auth/command"
	"upworkapi/internal/shared/contract"
	"upworkapi/internal/shared/contract/params"
)

type signUpEndpoint struct {
	params.AuthRouteParams
}

func NewSignUpEndpoint(routeParams params.AuthRouteParams) contract.Endpoint {
	return &signUpEndpoint{
		AuthRouteParams: routeParams,
	}
}

func (ep *signUpEndpoint) MapEndpoint() {
	ep.RootGroup.POST("/signup", ep.signUpHandler())
}

type signUpRequest struct {
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8"`
	ConfirmationPassword string `json:"confirmationPassword" validate:"required,eqfield=Password"`
}

func (ep *signUpEndpoint) signUpHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req signUpRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// Determine if the user already exists
		cmd := &command.GetSupabaseUserByEmail{Email: req.Email}
		_, err := mediatr.Send[*command.GetSupabaseUserByEmail, *supabase.User](ctx, cmd)
		if err == nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Email already exists")
		} else if !errors.Is(err, command.ErrUserNotFound) {
			ep.Logger.Error("Could not determine if user exists", "err", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Error checking user")
		}

		// Sign the user up
		user, err := ep.Supabase.Client.Auth.SignUp(c.Request().Context(), supabase.UserCredentials{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil && err.Error() != "Email rate limit exceeded" {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, user)
	}
}
