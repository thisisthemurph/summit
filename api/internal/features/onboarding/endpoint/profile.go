package endpoint

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"net/http"
	"upworkapi/internal/features/onboarding/command"
	"upworkapi/internal/shared/auth"
	sharedCmd "upworkapi/internal/shared/command"
	"upworkapi/internal/shared/contract/params"
	"upworkapi/internal/shared/model"
)

type OnboardingProfileEndpoint struct {
	params.OnboardingRouteParams
}

func NewOnboardingProfileEndpoint(routeParams params.OnboardingRouteParams) *OnboardingProfileEndpoint {
	return &OnboardingProfileEndpoint{
		OnboardingRouteParams: routeParams,
	}
}

func (ep *OnboardingProfileEndpoint) MapEndpoint() {
	ep.OnboardingGroup.GET("/profile", ep.getProfileData())
	ep.OnboardingGroup.POST("/profile", ep.upsertProfileData())
}

type GetProfileResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (ep *OnboardingProfileEndpoint) getProfileData() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := auth.GetAuthenticatedUser(c)

		query := &sharedCmd.GetUserByIDQuery{ID: user.ID}
		userRes, err := mediatr.Send[*sharedCmd.GetUserByIDQuery, *model.User](ctx, query)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, GetProfileResponse{
			FirstName: userRes.FirstName,
			LastName:  userRes.LastName,
		})
	}
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName" validate:"required,min=1,max=32"`
	LastName  string `json:"lastName" validate:"required,min=1,max=32"`
}

func (ep *OnboardingProfileEndpoint) upsertProfileData() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := auth.GetAuthenticatedUser(c)

		var req UpdateProfileRequest
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		cmd := command.NewUpdateProfileCommand(user.ID, req.FirstName, req.LastName)
		_, err := mediatr.Send[*command.UpdateProfileCommand, bool](ctx, cmd)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		query := &sharedCmd.GetUserByIDQuery{ID: user.ID}
		userResult, err := mediatr.Send[*sharedCmd.GetUserByIDQuery, *model.User](ctx, query)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		user.FirstName = userResult.FirstName
		user.LastName = userResult.LastName
		auth.SaveAuthenticatedUserInContext(c, user)
		return c.JSON(http.StatusOK, userResult)
	}
}
