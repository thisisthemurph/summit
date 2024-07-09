package application

import (
	"upworkapi/internal/shared/contract"
)

func (app *Application) MapEndpoints() {
	endpoints := app.Container.Get("routes").([]contract.Endpoint)
	for _, endpoint := range endpoints {
		endpoint.MapEndpoint()
	}

	//e := app.Container.Get("echo").(*echo.Echo)
	//e.Static("/public", "internal/app/shared/ui/static")
}
