package builder

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sarulabs/di"
	"log/slog"
	"os"
	"upworkapi/cmd/api/application"
	"upworkapi/internal/shared/contract"
	"upworkapi/internal/shared/contract/params"
	mw "upworkapi/internal/shared/middleware"
	"upworkapi/pkg/db"
	"upworkapi/pkg/supa"

	authEndpoint "upworkapi/internal/features/auth/endpoint"
	onboardingEndpoint "upworkapi/internal/features/onboarding/endpoint"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (b *ApplicationBuilder) AddCore() {
	err := godotenv.Load()
	check(err)

	logDep := di.Def{
		Name:  "log",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return b.Logger, nil
		},
	}

	configDep := di.Def{
		Name:  "config",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return application.NewConfig(os.Getenv), nil
		},
	}

	err = b.Services.Add(logDep)
	check(err)

	err = b.Services.Add(configDep)
	check(err)
}

func (b *ApplicationBuilder) AddInfrastructure() {
	err := addEcho(b.Services)
	check(err)

	supabaseDep := di.Def{
		Name:  "supabase",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get("config").(*application.Config)
			return supa.New(cfg.Supabase.URL, cfg.Supabase.Secret), nil
		},
	}

	dbDep := di.Def{
		Name:  "db",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cnf := ctn.Get("config").(*application.Config)
			database, err := db.Connect(cnf.Database)
			if err != nil {
				return nil, err
			}
			return database, nil
		},
	}

	check(b.Services.Add(supabaseDep))
	check(b.Services.Add(dbDep))
}

func (b *ApplicationBuilder) AddRoutes() {
	routesDep := di.Def{
		Name:  "routes",
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			authRouteParams := ctn.Get("auth_route_group").(*params.AuthRouteParams)
			onboardingRouteParams := ctn.Get("onboarding_route_group").(*params.OnboardingRouteParams)

			endpoints := []contract.Endpoint{
				authEndpoint.NewLoginEndpoint(*authRouteParams),
				authEndpoint.NewSignUpEndpoint(*authRouteParams),
				authEndpoint.NewLogOutEndpoint(*authRouteParams),
				onboardingEndpoint.NewOnboardingProfileEndpoint(*onboardingRouteParams),
			}

			return endpoints, nil
		},
	}

	err := b.Services.Add(routesDep)
	check(err)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func addEcho(container *di.Builder) error {
	deps := []di.Def{
		{
			Name:  "echo",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				logger := ctn.Get("log").(*slog.Logger)
				config := ctn.Get("config").(*application.Config)
				sb := ctn.Get("supabase").(*supa.Supabase)

				e := echo.New()
				e.Validator = &CustomValidator{validator: validator.New()}

				// Middleware
				authMiddleware := mw.NewAuthMiddleware(config, sb, logger)
				e.Use(middleware.Secure())
				e.Use(authMiddleware.WithUser)

				// CORS Middleware
				e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
					AllowOriginFunc: func(origin string) (bool, error) {
						return true, nil
					},
					AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
					AllowCredentials: true,
				}))
				return e, nil
			},
		},
		{
			Name:  "api_v1",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				e := ctn.Get("echo").(*echo.Echo)
				return e.Group("/api/v1"), nil
			},
		},
		{
			Name:  "auth_route_group",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				api := getApiGroup(ctn)
				logger := ctn.Get("log").(*slog.Logger)
				sb := ctn.Get("supabase").(*supa.Supabase)
				config := ctn.Get("config").(*application.Config)

				authRouteParams := &params.AuthRouteParams{
					RootGroup:     api.Group(""),
					SessionSecret: config.SessionSecret,
					Logger:        logger,
					Supabase:      sb,
				}

				return authRouteParams, nil
			},
		},
		{
			Name:  "onboarding_route_group",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				group := makeAuthGroup(ctn, "/onboarding")
				logger := ctn.Get("log").(*slog.Logger)

				onboardingRouteParams := &params.OnboardingRouteParams{
					OnboardingGroup: group,
					Logger:          logger,
				}

				return onboardingRouteParams, nil
			},
		},
	}

	for _, dep := range deps {
		if err := container.Add(dep); err != nil {
			return err
		}
	}
	return nil
}

func getApiGroup(ctn di.Container) *echo.Group {
	return ctn.Get("api_v1").(*echo.Group)
}

func makeAuthGroup(ctn di.Container, prefix string) *echo.Group {
	logger := ctn.Get("log").(*slog.Logger)
	sb := ctn.Get("supabase").(*supa.Supabase)
	cfg := ctn.Get("config").(*application.Config)
	baseGroup := getApiGroup(ctn)

	authMw := mw.NewAuthMiddleware(cfg, sb, logger)
	return baseGroup.Group(prefix, authMw.WithAuthenticatedUser)
}
