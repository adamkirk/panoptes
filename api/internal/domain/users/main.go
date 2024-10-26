package users

import (
	"time"

	"github.com/google/uuid"
)

type Encrypter interface {
	Encrypt(string) (string, error)
}

type AccessToken struct {
	ID string

	// Only set when the token is first generated
	Secret *string
	SecretHash string

	ExpireAt *time.Time

	User *User
}

type User struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
	PasswordHash string

	Roles []*Role
}

func (u *User) HasRole(search string) bool {
	for _, role := range u.Roles {
		if role.Name == search {
			return true
		}
	}

	return false
}

func (u *User) Can(actions []string) bool {
	if u.HasRole("superuser") {
		return true
	}

	return false
}

type Permission struct {
	ID uuid.UUID
	Name string
}

type Role struct {
	ID uuid.UUID
	Name string

	Permissions []*Permission
}