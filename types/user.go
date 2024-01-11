package types

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2/log"
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

func (params CreateUserParams) Validate() map[string]string {
	log.Info("enter Validate()")
	errors := map[string]string{}

	if len(params.FirstName) < minLenFirstname {
		errors["firstName"] = fmt.Sprintf("firstName should be atleast %d characters", minLenFirstname)
	}
	if len(params.LastName) < minLenLastname {
		errors["lastName"] = fmt.Sprintf("lastName should be atleast %d characters", minLenLastname)
	}
	if len(params.Email) < minLenPassword {
		errors["password"] = fmt.Sprintf("password should be atleast %d characters", minLenPassword)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	log.Info("exit Validate()")
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	log.Info("enter NewUserFromParams")
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

/* type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (params UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(params.FirstName) > 0 {
		m["firstName"] = params.FirstName
	}

	if len(params.LastName) > 0 {
		m["lastName"] = params.LastName
	}
	return m
} */
