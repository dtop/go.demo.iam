package models

import (
	"database/sql"
	"regexp"

	"crypto/sha1"
	"fmt"
	"io"

	"log"

	"github.com/dtop/go.demo.iam/iam/wrappers"
	"github.com/dtop/go.ginject"
)

type (

	// UserLogin is the user login handler
	UserLogin struct {
		UserID    string
		SessionID string `form:"sess" binding:"required"`
		Email     string `form:"eml" binding:"required"`
		Password  string `form:"pwd" binding:"required"`

		Db     *wrappers.MySQL `inject:"db"`
		errors []Err
	}

	Err struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

// NewUserLogin returns a brand new user login model
func NewUserLogin(deps ginject.Injector) *UserLogin {

	ulog := &UserLogin{errors: make([]Err, 0)}
	deps.Apply(ulog)

	return ulog
}

// CheckLogin checks if eml/pwd is valid and assigns the userID if so
func (ul *UserLogin) CheckLogin() []Err {

	if ul.SessionID == "" {
		ul.addError("SessionID was missing, please start over", "none")
	}

	if ul.Email == "" || !isEmail(ul.Email) {
		ul.addError("Please enter a valid email address", "eml")
	}

	if len(ul.Password) < 6 {
		ul.addError("Please enter a valid password", "pwd")
	}

	if len(ul.errors) > 0 {
		return ul.errors
	}

	result, err := ul.Db.Query(func(db *sql.DB) (*sql.Rows, error) {

		query := "SELECT id, password FROM user_data WHERE email LIKE ?"
		return db.Query(query, ul.Email)
	})

	if err != nil {
		ul.addError("server error occured, please try later", "none")
		return ul.errors
	}

	if !result.Next() {
		log.Println("unknown record")
		ul.addError("could not find this username/password combination", "none")
		return ul.errors
	}

	var (
		ident string
		pass  string
	)

	if err := result.Scan(&ident, &pass); err != nil {
		ul.addError("server error occured, please try later", "none")
		return ul.errors
	}

	if hashPassword(ul.Password) != pass {
		log.Println(hashPassword(ul.Password), " does not match ", pass)
		ul.addError("could not find this username/password combination", "none")
		return ul.errors
	}

	ul.UserID = ident
	return ul.errors
}

func (ul *UserLogin) addError(text, field string) {

	ul.errors = append(ul.errors, Err{Field: field, Message: text})
}

// Adopted from StackOverflow =)
func isEmail(mail string) bool {

	pttrn := "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	rxEmail := regexp.MustCompile(pttrn)
	return rxEmail.MatchString(mail)
}

func hashPassword(pwd string) string {

	h := sha1.New()
	io.WriteString(h, pwd)
	return fmt.Sprintf("%x", h.Sum(nil))
}
