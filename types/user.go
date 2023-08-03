package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (p CreateUserParams) Validate() []string {
	var errors []string
	if len(p.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("firstname length should be at least %d characters", minFirstNameLen))
	}
	if len(p.LastName) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen))
	}
	if len(p.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password length should be at least %d characters", minPasswordLen))
	}
	if len(p.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("firstname length should be at leasdt %d characters", minFirstNameLen))
	}
	if !isEmailValid(p.Email) {
		errors = append(errors, fmt.Sprintf("email is invalid"))
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegExp := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegExp.MatchString(e)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"LastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	fmt.Println(encpw)
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
