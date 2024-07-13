package command

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type UpdateProfileCommand struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
}

func NewUpdateProfileCommand(userID uuid.UUID, firstName, lastName string) *UpdateProfileCommand {
	return &UpdateProfileCommand{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
	}
}

type UpdateProfileHandler struct {
	DB *sql.DB
}

func (h *UpdateProfileHandler) Handle(ctx context.Context, cmd *UpdateProfileCommand) (bool, error) {
	stmt := `
		INSERT INTO user_profiles (id, first_name, last_name) VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET first_name = $2,
		    last_name = $3;`

	_, err := h.DB.ExecContext(ctx, stmt, cmd.UserID, cmd.FirstName, cmd.LastName)
	if err != nil {
		return false, err
	}
	return true, nil
}
