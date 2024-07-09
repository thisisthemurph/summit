package builder

import (
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
	"log/slog"
	"os"
	"upworkapi/cmd/api/application"
)

type ApplicationBuilder struct {
	Services *di.Builder
	Logger   *slog.Logger
}

func NewApplicationBuilder() *ApplicationBuilder {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	builder, err := di.NewBuilder()
	if err != nil {
		panic(err)
	}

	return &ApplicationBuilder{
		Services: builder,
		Logger:   logger,
	}
}

func (b *ApplicationBuilder) Build() *application.Application {
	container := b.Services.Build()

	e := container.Get("echo").(*echo.Echo)
	cfg := container.Get("config").(*application.Config)
	return application.New(container, e, b.Logger, cfg)
}
