package builder

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sarulabs/di"
	"log/slog"
	"os"
	"upworkapi/cmd/api/application"
	authEndpoint "upworkapi/internal/features/auth/endpoint"
	"upworkapi/internal/shared/contract"
	"upworkapi/internal/shared/contract/params"
	"upworkapi/pkg/db"
	"upworkapi/pkg/supa"
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

			endpoints := []contract.Endpoint{
				authEndpoint.NewLoginEndpoint(*authRouteParams),
				authEndpoint.NewSignUpEndpoint(*authRouteParams),
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
				//logger := ctn.Get("log").(*slog.Logger)
				//cfg := ctn.Get("config").(*application.Config)

				e := echo.New()
				e.Validator = &CustomValidator{validator: validator.New()}

				// Middleware
				e.Use(middleware.Secure())

				// CORS Middleware
				e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
					//AllowOrigins: []string{"http://localhost:5173", "http://localhost:8080"},
					AllowOriginFunc: func(origin string) (bool, error) {
						fmt.Println(origin)
						return true, nil
					},
					AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
				}))

				//e.Static("/public", "internal/app/shared/ui/static")
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

				authRouteParams := &params.AuthRouteParams{
					LoginGroup:  api.Group("/login"),
					SignupGroup: api.Group("/signup"),
					LogoutGroup: api.Group("/logout"),
					Logger:      logger,
					Supabase:    sb,
				}

				return authRouteParams, nil
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
