package command

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"upworkapi/internal/shared/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type GetUserByIDQuery struct {
	ID uuid.UUID
}

type GetUserByIDHandler struct {
	DB *sql.DB
}

func (h *GetUserByIDHandler) Handle(ctx context.Context, req *GetUserByIDQuery) (*model.User, error) {
	query := `
		SELECT u.id, p.first_name, p.last_name, u.email
		FROM auth.users u
		LEFT JOIN public.user_profiles p ON u.id = p.id
		WHERE u.id = $1
		LIMIT 1;`

	var (
		id    uuid.UUID
		fname sql.NullString
		lname sql.NullString
		email string
	)

	err := h.DB.QueryRowContext(ctx, query, req.ID).Scan(&id, &fname, &lname, &email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &model.User{
		ID:        id,
		FirstName: fname.String,
		LastName:  lname.String,
		Email:     email,
	}, nil
}
