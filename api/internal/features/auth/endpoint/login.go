package endpoint

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/nedpals/supabase-go"
	"net/http"
	"upworkapi/internal/shared/auth"
	"upworkapi/internal/shared/command"
	"upworkapi/internal/shared/contract"
	"upworkapi/internal/shared/contract/params"
	"upworkapi/internal/shared/model"
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
	ep.RootGroup.POST("/login", ep.loginHandler())
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

		// Sign the user in using Supabase

		credentials := supabase.UserCredentials{
			Email:    req.Email,
			Password: req.Password,
		}
		authDetails, err := ep.Supabase.Client.Auth.SignIn(ctx, credentials)

		if err != nil {
			ep.Logger.Warn("login failed", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidCredentials.Error())
		}

		// Set the auth session cookie

		if err := auth.SetAuthSession(c, authDetails.AccessToken, ep.SessionSecret); err != nil {
			ep.Logger.Error("could not set auth session", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		// Get user information from the database

		userID, _ := uuid.Parse(authDetails.User.ID)
		cmd := &command.GetUserByIDQuery{ID: userID}
		user, err := mediatr.Send[*command.GetUserByIDQuery, *model.User](ctx, cmd)
		if err != nil {
			ep.Logger.Warn("could not get user", "error", err)
			if errors.Is(err, command.ErrUserNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "user not found")
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		authenticatedUser := auth.AuthenticatedUser{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		}
		return c.JSON(http.StatusOK, authenticatedUser)
	}
}
