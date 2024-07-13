package application

import (
	"database/sql"
	"errors"
	"github.com/mehdihadeli/go-mediatr"
	"upworkapi/internal/features/onboarding/command"
	"upworkapi/internal/shared/model"

	shared "upworkapi/internal/shared/command"
)

func (app *Application) ConfigureMediator() error {
	db, ok := app.Container.Get("db").(*sql.DB)
	if !ok {
		return errors.New("db connection error")
	}

	getUserByIdHandler := &shared.GetUserByIDHandler{DB: db}
	err := mediatr.RegisterRequestHandler[*shared.GetUserByIDQuery, *model.User](getUserByIdHandler)
	if err != nil {
		return err
	}

	getUserByEmailHandler := &shared.GetUserByEmailHandler{DB: db}
	err = mediatr.RegisterRequestHandler[*shared.GetUserByEmail, *model.User](getUserByEmailHandler)
	if err != nil {
		return err
	}

	updateProfileHandler := &command.UpdateProfileHandler{DB: db}
	err = mediatr.RegisterRequestHandler[*command.UpdateProfileCommand, bool](updateProfileHandler)
	if err != nil {
		return err
	}

	return nil
}
