package main

import "upworkapi/cmd/api/builder"

func main() {
	b := builder.NewApplicationBuilder()
	b.AddCore()
	b.AddInfrastructure()
	b.AddRoutes()

	app := b.Build()
	if err := app.ConfigureMediator(); err != nil {
		panic(err)
	}
	app.MapEndpoints()
	app.Run()
}
