package models

import (
	"fmt"
	"regexp"

	"github.com/revel/revel"
)

const (
	// User Model Consts
	UserMinUsernameLength int = 4
	UserMaxUsernameLength int = 20
	UserMaxNameLength     int = 40
	UserMinPasswordLength int = 5
	UserMaxPasswordLength int = 50
	UserMaxEmailLength    int = 200

	// AccountType Model Consts
	UserAccountTypeAdmin string = "ADMIN"
	UserAccountTypeUser  string = "USER"
)

var (
	userRegex = regexp.MustCompile("^\\w*$")
)

type User struct {
	UserId         int
	Name           string
	Username       string
	Email          string
	AccountType    string
	Password       string // This isn't persisted, but is here for databinding from the front-end
	HashedPassword []byte
}

type AccountType struct {
	AccountTypeId int
	Name          string
	Description   string
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{UserMaxUsernameLength},
		revel.MinSize{UserMinUsernameLength},
		revel.Match{userRegex},
	)

	ValidatePassword(v, user.Password).
		Key("user.Password")

	v.Check(user.Name,
		revel.Required{},
		revel.MaxSize{UserMaxNameLength},
	)

	v.Check(user.Email,
		revel.Required{},
		revel.MaxSize{UserMaxEmailLength},
	)

	v.Email(user.Email)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{UserMaxPasswordLength},
		revel.MinSize{UserMinPasswordLength},
	)
}
