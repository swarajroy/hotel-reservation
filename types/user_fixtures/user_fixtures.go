package userfixtures

import (
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/swarajroy/hotel-reservation/types"
	"golang.org/x/crypto/bcrypt"
)

const (
	BCRYPT_COST = 12
)

func Next() (*types.User, error) {
	u := &types.User{}
	err := faker.FakeData(u)
	if err != nil {
		return nil, fmt.Errorf("error")
	}
	return u, nil
}

func NextWith(fn, ln, email, password string, isAdmin bool) (*types.User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_COST)
	if err != nil {
		return nil, fmt.Errorf("error generating encrypted password")
	}
	return &types.User{
		FirstName:         fn,
		LastName:          ln,
		Email:             email,
		EncryptedPassword: string(encryptedPassword),
		IsAdmin:           isAdmin,
	}, nil
}
