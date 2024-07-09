package application

import (
	"database/sql"
	"errors"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/nedpals/supabase-go"
	"upworkapi/internal/features/auth/command"
)

func (app *Application) ConfigureMediator() error {
	db, ok := app.Container.Get("db").(*sql.DB)
	if !ok {
		return errors.New("db connection error")
	}

	getUserByEmailHandler := &command.GetSupabaseUserByEmailHandler{DB: db}
	err := mediatr.RegisterRequestHandler[*command.GetSupabaseUserByEmail, *supabase.User](getUserByEmailHandler)
	if err != nil {
		return err
	}

	return nil
}
