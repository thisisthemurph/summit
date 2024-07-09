package command

import (
	"context"
	"database/sql"
	"errors"
	"github.com/nedpals/supabase-go"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type GetSupabaseUserByEmail struct {
	Email string
}

type GetSupabaseUserByEmailHandler struct {
	DB *sql.DB
}

type UserMetadata struct {
	EmailVerified bool `json:"email_verified"`
	PhoneVerified bool `json:"phone_verified"`
}

func (h *GetSupabaseUserByEmailHandler) Handle(ctx context.Context, req *GetSupabaseUserByEmail) (*supabase.User, error) {
	query := `
		SELECT id, email
		FROM auth.users
		WHERE email = $1
		LIMIT 1;`

	var user supabase.User
	err := h.DB.QueryRowContext(ctx, query, req.Email).Scan(
		&user.ID, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
