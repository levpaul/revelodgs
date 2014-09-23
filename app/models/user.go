package models

import (
	"fmt"
	"regexp"

	"github.com/revel/revel"
)

const (
	// User Model Consts
	UserMaxUsernameLength int = 150
	UserMaxNameLength     int = 20
	UserMaxPasswordLength int = 40

	// AccountType Model Consts
	AccountTypeAdminName string = "ADMIN"
	AccountTypeAdminId   int    = 1
	AccountTypeUserName  string = "USER"
	AccountTypeUserId    int    = 2
)

type User struct {
	UserId         int
	Name           string
	Username       string
	Email          string
	AccountType    int
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

var userRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
		revel.Match{userRegex},
	)

	ValidatePassword(v, user.Password).
		Key("user.Password")

	v.Check(user.Name,
		revel.Required{},
		revel.MaxSize{100},
	)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{5},
	)
}
