package auth

import "github.com/google/uuid"

type AuthenticatedUser struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	AccessToken string    `json:"-"`
	LoggedIn    bool      `json:"-"`
}

// ProfileComplete returns false if any of the required profile fields are not set, otherwise true.
func (u AuthenticatedUser) ProfileComplete() bool {
	if u.FirstName == "" || u.LastName == "" {
		return false
	}

	return true
}
