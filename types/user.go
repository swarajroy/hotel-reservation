package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	BCRYPT_COST     = 12
	minLenFirstname = 2
	minLenLastname  = 2
	minLenPassword  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() []string {
	errors := []string{}

	if len(params.FirstName) < minLenFirstname {
		errors = append(errors, fmt.Sprintf("firstName should be atleast %d characters", minLenFirstname))
	}
	if len(params.LastName) < minLenLastname {
		errors = append(errors, fmt.Sprintf("lastName should be atleast %d characters", minLenLastname))
	}
	if len(params.Email) < minLenPassword {
		errors = append(errors, fmt.Sprintf("password should be atleast %d characters", minLenPassword))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, fmt.Sprintf("email is invalid"))
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), BCRYPT_COST)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}
